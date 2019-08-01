package watcher

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	"github.com/rjeczalik/notify"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/util/sync"
)

const eventBuffer = 4096 // How many events and updates we can buffer before watching is interrupted
var listenEvents = []notify.Event{notify.InCreate, notify.InDelete, notify.InDeleteSelf, notify.InCloseWrite}

var eventMap = map[notify.Event]Event{
	notify.InCreate:     EventCreate,
	notify.InDelete:     EventDelete,
	notify.InCloseWrite: EventModify,
}

// combinedEvent describes multiple events that should be concatenated into a single event
type combinedEvent struct {
	input  []byte // input is a slice of events to match (in bytes, it speeds up the comparison)
	output Event  // output is the resulting event that should be returned
}

// combinedEvents describes the event combinations to concatenate,
// this is iterated in order, so the longest matches should be first
var combinedEvents = []combinedEvent{
	// DELETE + CREATE + MODIFY => MODIFY
	{Events{EventDelete, EventCreate, EventModify}.Bytes(), EventModify},
	// CREATE + MODIFY => CREATE
	{Events{EventCreate, EventModify}.Bytes(), EventCreate},
	// CREATE + DELETE => NONE
	{Events{EventCreate, EventDelete}.Bytes(), EventNone},
}

// Suppress duplicate events registered in this map. E.g. directory deletion
// fires two DELETE events, one for the parent and one for the deleted directory itself
var suppressDuplicates = map[Event]bool{
	EventCreate: true,
	EventDelete: true,
}

type eventStream chan notify.EventInfo
type FileUpdateStream chan *FileUpdate
type watches []string

// Options specifies options for the FileWatcher
type Options struct {
	// ExcludeDirs specifies what directories to not watch
	ExcludeDirs []string
	// BatchTimeout specifies the duration to wait after last event before dispatching grouped inotify events
	BatchTimeout time.Duration
	// ValidExtensions specifies what file extensions to look at
	ValidExtensions []string
}

// DefaultOptions returns the default options
func DefaultOptions() Options {
	return Options{
		ExcludeDirs:     []string{".git"},
		BatchTimeout:    1 * time.Second,
		ValidExtensions: []string{".yaml", ".yml", ".json"},
	}
}

// NewFileWatcher returns a list of files in the watched directory in
// addition to the generated FileWatcher, it can be used to populate
// MappedRawStorage fileMappings
func NewFileWatcher(dir string) (w *FileWatcher, files []string, err error) {
	return NewFileWatcherWithOptions(dir, DefaultOptions())
}

// NewFileWatcher returns a list of files in the watched directory in
// addition to the generated FileWatcher, it can be used to populate
// MappedRawStorage fileMappings
func NewFileWatcherWithOptions(dir string, opts Options) (w *FileWatcher, files []string, err error) {
	w = &FileWatcher{
		dir:     dir,
		events:  make(eventStream, eventBuffer),
		updates: make(FileUpdateStream, eventBuffer),
		batcher: sync.NewBatchWriter(opts.BatchTimeout),
		opts:    opts,
	}

	if err = w.start(&files); err != nil {
		notify.Stop(w.events)
	} else {
		w.monitor = sync.RunMonitor(w.monitorFunc)
		w.dispatcher = sync.RunMonitor(w.dispatchFunc)
	}

	return
}

// FileWatcher recursively monitors changes in files in the given directory
// and sends out events based on their state changes. Only files conforming
// to validSuffix are monitored. The FileWatcher can be suspended for a single
// event at a time to eliminate updates by WatchStorage causing a loop.
type FileWatcher struct {
	dir          string
	events       eventStream
	updates      FileUpdateStream
	watches      watches
	suspendEvent Event
	monitor      *sync.Monitor
	dispatcher   *sync.Monitor
	opts         Options
	// the batcher is used for properly sending many concurrent inotify events
	// as a group, after a specified timeout. This fixes the issue of one single
	// file operation being registered as many different inotify events
	batcher *sync.BatchWriter
}

// TODO: Add a recursive watch to root instead of one for every path
func (w *FileWatcher) addWatch(path string) (err error) {
	log.Tracef("FileWatcher: Adding watch for %q", path)
	if err = notify.Watch(path, w.events, listenEvents...); err == nil {
		w.watches = append(w.watches, path)
	}

	return
}

func (w *FileWatcher) hasWatch(path string) bool {
	for _, watch := range w.watches {
		if watch == path {
			log.Tracef("FileWatcher: Watch found for %q", path)
			return true
		}
	}

	log.Tracef("FileWatcher: No watch found for %q", path)
	return false
}

func (w *FileWatcher) clear() {
	log.Tracef("FileWatcher: Clearing all watches")
	notify.Stop(w.events)
	w.watches = w.watches[:0]
}

// start discovers all subdirectories and adds paths to
// notify before starting the monitoring goroutine
func (w *FileWatcher) start(files *[]string) error {
	return filepath.Walk(w.dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				for _, dir := range w.opts.ExcludeDirs {
					if info.Name() == dir {
						return filepath.SkipDir // Skip excluded directories
					}
				}

				return w.addWatch(path)
			}

			if files != nil {
				// Only include files with a valid suffix
				if w.validSuffix(info.Name()) {
					*files = append(*files, path)
				}
			}

			return nil
		})
}

func (w *FileWatcher) monitorFunc() {
	log.Debug("FileWatcher: Monitoring thread started")
	defer log.Debug("FileWatcher: Monitoring thread stopped")
	defer close(w.updates) // Close the update stream after the FileWatcher has stopped

	for {
		event, ok := <-w.events
		if !ok {
			return
		}

		updateEvent := convertEvent(event.Event())
		if updateEvent == w.suspendEvent {
			w.suspendEvent = 0
			log.Debugf("FileWatcher: Skipping suspended event %s for path: %q", updateEvent, event.Path())
			continue // Skip the suspended event
		}

		// Suppress successive duplicate events registered in suppressDuplicates
		if suppressEvent(event.Path(), updateEvent) {
			log.Debugf("FileWatcher: Skipping suppressed event %s for path: %q", updateEvent, event.Path())
			continue // Skip the suppressed event
		}

		// Directory bypass for FileWatcher registration
		// The FileWatcher registration/deletion needs to be as fast as
		// possible, bypass the batcher when dealing with directories
		if w.handleDirEvent(event.Path(), updateEvent) {
			continue // The event path matched a directory, skip file processing for the event
		}

		// Get any events registered for the specific file, and append the specified event
		var eventList Events
		if val, ok := w.batcher.Load(event.Path()); ok {
			eventList = val.(Events)
		}

		eventList = append(eventList, updateEvent)

		// Register the event in the map, and dispatch all the events at once after the timeout
		w.batcher.Store(event.Path(), eventList)
		log.Debugf("FileWatcher: Registered inotify events %v for path %q", eventList, event.Path())
	}
}

func (w *FileWatcher) dispatchFunc() {
	log.Debug("FileWatcher: Dispatch thread started")
	defer log.Debug("FileWatcher: Dispatch thread stopped")

	for {
		// Wait until we have a batch dispatched to us
		ok := w.batcher.ProcessBatch(func(key, val interface{}) bool {
			filePath := key.(string)

			// Concatenate all known events, and dispatch them to be handled one by one
			for _, event := range concatenateEvents(val.(Events)) {
				w.handleEvent(filePath, event)
			}

			// Continue traversing the map
			return true
		})
		if !ok {
			return // The BatchWriter channel is closed, stop processing
		}

		log.Debug("FileWatcher: Dispatched events batch and reset the events cache")
	}
}

func (w *FileWatcher) handleEvent(filePath string, event Event) {
	switch event {
	case EventCreate, EventDelete, EventModify: // Ignore EventNone
		// only care about valid files
		if !w.validSuffix(filePath) {
			return
		}

		log.Debugf("FileWatcher: Sending update: %s -> %q", event, filePath)
		w.updates <- &FileUpdate{
			Event: event,
			Path:  filePath,
		}
	}
}

func (w *FileWatcher) handleDirEvent(filePath string, event Event) (dir bool) {
	switch event {
	case EventCreate:
		fi, err := os.Stat(filePath)
		if err != nil {
			log.Errorf("FileWatcher: Failed to stat %q: %v", filePath, err)
			return
		}

		if fi.IsDir() {
			if err := w.addWatch(filePath); err != nil {
				log.Errorf("FileWatcher: Failed to add %q: %v", filePath, err)
			}

			dir = true
		}
	case EventDelete:
		if w.hasWatch(filePath) {
			w.clear()
			if err := w.start(nil); err != nil {
				log.Errorf("FileWatcher: Failed to re-initialize watches for %q", w.dir)
			}

			dir = true
		}
	}

	return
}

// GetFileUpdateStream gets the channel with FileUpdates
func (w *FileWatcher) GetFileUpdateStream() FileUpdateStream {
	return w.updates
}

// Close closes active underlying resources
func (w *FileWatcher) Close() {
	notify.Stop(w.events)
	w.batcher.Close()
	close(w.events) // Close the event stream
	w.monitor.Wait()
	w.dispatcher.Wait()
}

// Suspend enables a one-time suspend of the given event,
// the FileWatcher will skip the given event once
func (w *FileWatcher) Suspend(updateEvent Event) {
	w.suspendEvent = updateEvent
}

// validSuffix is used to filter out all unsupported
// files based on the extensions given
func (w *FileWatcher) validSuffix(path string) bool {
	for _, suffix := range w.opts.ValidExtensions {
		if filepath.Ext(path) == suffix {
			return true
		}
	}

	return false
}

func convertEvent(event notify.Event) Event {
	if updateEvent, ok := eventMap[event]; ok {
		return updateEvent
	}

	return EventNone
}

// concatenateEvents takes in a slice of events and concatenates
// all events possible based on combinedEvents
func concatenateEvents(events Events) Events {
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
			concatenated := append(Events{combinedEvent.output}, events[len(combinedEvent.input):]...)
			log.Tracef("FileWatcher: Concatenated events: %v -> %v", events, concatenated)
			return concatenateEvents(concatenated)
		}
	}

	return events
}

var suppressCache struct {
	event Event
	path  string
}

// suppressEvent returns true it it's called twice
// in a row with the same known event and path
func suppressEvent(path string, event Event) (s bool) {
	if _, ok := suppressDuplicates[event]; ok {
		if suppressCache.event == event && suppressCache.path == path {
			s = true
		}
	}

	suppressCache.event = event
	suppressCache.path = path
	return
}
