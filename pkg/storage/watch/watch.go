package watch

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	"github.com/rjeczalik/notify"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/storage"
	"github.com/weaveworks/ignite/pkg/storage/watch/update"
	"github.com/weaveworks/ignite/pkg/util"
)

const eventBuffer = 4096                 // How many events and updates we can buffer before watching is interrupted
const dispatchDuration = 1 * time.Second // Duration to wait after last event before dispatching grouped inotify events
var excludeDirs = []string{".git"}
var listenEvents = []notify.Event{notify.InCreate, notify.InDelete, notify.InDeleteSelf, notify.InCloseWrite}

var eventMap = map[notify.Event]update.Event{
	notify.InCreate:     update.EventCreate,
	notify.InDelete:     update.EventDelete,
	notify.InCloseWrite: update.EventModify,
}

// combinedEvent describes multiple events that should be concatenated into a single event
type combinedEvent struct {
	input  []byte       // input is a slice of events to match (in bytes, it speeds up the comparison)
	output update.Event // output is the resulting event that should be returned
}

// combinedEvents describes the event combinations to concatenate,
// this is iterated in order, so the longest matches should be first
var combinedEvents = []combinedEvent{
	// DELETE + CREATE + MODIFY => MODIFY
	{update.Events{update.EventDelete, update.EventCreate, update.EventModify}.Bytes(), update.EventModify},
	// CREATE + MODIFY => CREATE
	{update.Events{update.EventCreate, update.EventModify}.Bytes(), update.EventCreate},
	// CREATE + DELETE => NONE
	{update.Events{update.EventCreate, update.EventDelete}.Bytes(), update.EventNone},
}

type eventStream chan notify.EventInfo
type UpdateStream chan *update.FileUpdate
type watches []string

// watcher recursively monitors changes in files in the given directory
// and sends out events based on their state changes. Only files conforming
// to validSuffix are monitored. The watcher can be suspended for a single
// event at a time to eliminate updates by WatchStorage causing a loop.
type watcher struct {
	dir          string
	events       eventStream
	updates      UpdateStream
	watches      watches
	suspendEvent update.Event
	monitor      *util.Monitor
	dispatcher   *util.Monitor
	// the batcher is used for properly sending many concurrent inotify events
	// as a group, after a specified timeout. This fixes the issue of one single
	// file operation being registered as many different inotify events
	batcher *Batcher
}

func (w *watcher) addWatch(path string) (err error) {
	log.Tracef("Watcher: adding watch for %q", path)
	if err = notify.Watch(path, w.events, listenEvents...); err == nil {
		w.watches = append(w.watches, path)
	}

	return
}

func (w *watcher) hasWatch(path string) bool {
	for _, watch := range w.watches {
		if watch == path {
			log.Tracef("Watcher: watch found for %q", path)
			return true
		}
	}

	log.Tracef("Watcher: no watch found for %q", path)
	return false
}

func (w *watcher) clear() {
	log.Tracef("Watcher: clearing all watches")
	notify.Stop(w.events)
	w.watches = w.watches[:0]
}

// newWatcher returns a list of files in the watched directory in
// addition to the generated watcher, it can be used to populate
// MappedRawStorage fileMappings
func newWatcher(dir string) (w *watcher, files []string, err error) {
	w = &watcher{
		dir:     dir,
		events:  make(eventStream, eventBuffer),
		updates: make(UpdateStream, eventBuffer),
		batcher: NewBatcher(dispatchDuration),
	}

	if err = w.start(&files); err != nil {
		notify.Stop(w.events)
	} else {
		w.monitor = util.RunMonitor(w.monitorFunc)
		w.dispatcher = util.RunMonitor(w.dispatchFunc)
	}

	return
}

// start discovers all subdirectories and adds paths to
// notify before starting the monitoring goroutine
func (w *watcher) start(files *[]string) error {
	return filepath.Walk(w.dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				for _, dir := range excludeDirs {
					if info.Name() == dir {
						return filepath.SkipDir // Skip excluded directories
					}
				}

				return w.addWatch(path)
			}

			if files != nil {
				// Only include files with a valid suffix
				if validSuffix(info.Name()) {
					*files = append(*files, path)
				}
			}

			return nil
		})
}

func (w *watcher) monitorFunc() {
	log.Debug("Watcher: monitoring thread started")
	defer log.Debug("Watcher: monitoring thread stopped")
	defer close(w.updates) // Close the update stream after the watcher has stopped

	for {
		event, ok := <-w.events
		if !ok {
			return
		}

		updateEvent := convertEvent(event.Event())
		if updateEvent == w.suspendEvent {
			w.suspendEvent = 0
			log.Debugf("Watcher: skipping suspended event %s for path: %s", updateEvent, event.Path())
			continue // Skip the suspended event
		}

		// Get any events registered for the specific file, and append the specified event
		var eventList update.Events
		val, ok := w.batcher.Load(event.Path())
		if ok {
			eventList = val.(update.Events)
		}
		eventList = append(eventList, updateEvent)

		// Register the event in the map, and dispatch all the events at once after the timeout
		w.batcher.Store(event.Path(), eventList)
		log.Debugf("Watcher: Registered inotify events %v for path %s", eventList, event.Path())
	}
}

func (w *watcher) dispatchFunc() {
	log.Debug("Watcher: dispatch thread started")
	defer log.Debug("Watcher: dispatch thread stopped")

	for {
		// wait until we have a batch dispatched to us
		w.batcher.ProcessBatch(func(key, val interface{}) bool {
			filePath := key.(string)

			// Concatenate all known events, and dispatch them to be handled one by one
			for _, event := range concatenateEvents(val.(update.Events)) {
				w.handleEvent(filePath, event)
			}

			// continue traversing the map
			return true
		})
		log.Debug("Watcher: dispatched events batch and reset the events cache")
	}
}

func (w *watcher) handleEvent(filePath string, event update.Event) {
	switch event {
	case update.EventCreate:
		fi, err := os.Stat(filePath)
		if err != nil {
			log.Errorf("Watcher: failed to stat %q: %v", filePath, err)
			return
		}

		if fi.IsDir() {
			if err := w.addWatch(filePath); err != nil {
				log.Errorf("Watcher: failed to add %q: %v", filePath, err)
			}

			return
		}

		fallthrough
	case update.EventDelete, update.EventModify:
		if event == update.EventDelete && w.hasWatch(filePath) {
			w.clear()
			if err := w.start(nil); err != nil {
				log.Errorf("Watcher: Failed to re-initialize watches for %q", w.dir)
			}

			return
		}

		// only care about valid files
		if !validSuffix(filePath) {
			return
		}

		log.Debugf("Watcher: Sending update: %s -> %q", event, filePath)
		w.updates <- &update.FileUpdate{
			Event: event,
			Path:  filePath,
		}
	}
}

// TODO: This watcher doesn't handle multiple operations on the same file well
// DELETE+CREATE+MODIFY => MODIFY
// CREATE+MODIFY => CREATE
// Fix this by caching the operations on the same file, and one second after all operations
// have been "written"; go through the changes and interpret the combinations of events properly
// This maybe will allow us to remove the "suspend" functionality? I don't know yet

func (w *watcher) close() {
	notify.Stop(w.events)
	w.batcher.Close()
	close(w.events) // Close the event stream
	w.monitor.Wait()
	w.dispatcher.Wait()
}

// This enables a one-time suspend of the given event,
// the watcher will skip the given event once
func (w *watcher) suspend(updateEvent update.Event) {
	w.suspendEvent = updateEvent
}

func convertEvent(event notify.Event) update.Event {
	if updateEvent, ok := eventMap[event]; ok {
		return updateEvent
	}

	return update.EventNone
}

// validSuffix is used to filter out all unsupported
// files based on the extensions in storage.Formats
func validSuffix(path string) bool {
	for suffix := range storage.Formats {
		if filepath.Ext(path) == suffix {
			return true
		}
	}

	return false
}

// concatenateEvents takes in a slice of events and concatenates
// all events possible based on combinedEvents
func concatenateEvents(events update.Events) update.Events {
	if len(events) < 2 {
		return events // Quick return for 0 or 1 event
	}

	for _, combinedEvent := range combinedEvents {
		if len(combinedEvent.input) > len(events) {
			continue // The combined event's match is too long
		}

		// Test if the prefix of the given events matches combinedEvent.input
		if bytes.Equal(events.Bytes()[:len(combinedEvent.input)], combinedEvent.input) {
			// If so, replace combinedEvent.input prefix in events with combinedEvent.output and recurse
			concatenated := append(update.Events{combinedEvent.output}, events[len(combinedEvent.input):]...)
			log.Tracef("Watcher: concatenated events: %v -> %v", events, concatenated)
			return concatenateEvents(concatenated)
		}
	}

	return events
}
