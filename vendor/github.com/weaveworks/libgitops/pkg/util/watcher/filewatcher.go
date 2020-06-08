package watcher

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/rjeczalik/notify"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/libgitops/pkg/util/sync"
	"golang.org/x/sys/unix"
)

const eventBuffer = 4096 // How many events and updates we can buffer before watching is interrupted
var listenEvents = []notify.Event{notify.InDelete, notify.InCloseWrite, notify.InMovedFrom, notify.InMovedTo}

var eventMap = map[notify.Event]FileEvent{
	notify.InDelete:     FileEventDelete,
	notify.InCloseWrite: FileEventModify,
}

// combinedEvent describes multiple events that should be concatenated into a single event
type combinedEvent struct {
	input  []notify.Event // input is a slice of events to match (in bytes, it speeds up the comparison)
	output int            // output is the event's index that should be returned, negative values equal nil
}

func (c *combinedEvent) match(events notifyEvents) (notify.EventInfo, bool) {
	if len(c.input) > len(events) {
		return nil, false // Not enough events, cannot match
	}

	for i := 0; i < len(c.input); i++ {
		if events[i].Event() != c.input[i] {
			return nil, false
		}
	}

	if c.output > 0 {
		return events[c.output], true
	}

	return nil, true
}

// combinedEvents describes the event combinations to concatenate,
// this is iterated in order, so the longest matches should be first
var combinedEvents = []combinedEvent{
	// DELETE + MODIFY => MODIFY
	{[]notify.Event{notify.InDelete, notify.InCloseWrite}, 1},
	// MODIFY + DELETE => NONE
	{[]notify.Event{notify.InCloseWrite, notify.InDelete}, -1},
}

type notifyEvents []notify.EventInfo
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

	log.Tracef("FileWatcher: Starting recursive watch for %q", dir)
	if err = notify.Watch(path.Join(dir, "..."), w.events, listenEvents...); err != nil {
		notify.Stop(w.events)
	} else if files, err = w.getFiles(); err == nil {
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
	suspendEvent FileEvent
	monitor      *sync.Monitor
	dispatcher   *sync.Monitor
	opts         Options
	// the batcher is used for properly sending many concurrent inotify events
	// as a group, after a specified timeout. This fixes the issue of one single
	// file operation being registered as many different inotify events
	batcher *sync.BatchWriter
}

// getFiles discovers all subdirectories and
// returns a list of valid files in them
func (w *FileWatcher) getFiles() (files []string, err error) {
	err = filepath.Walk(w.dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				// Only include valid files
				if w.validFile(path) {
					files = append(files, path)
				}
			}

			return nil
		})

	return
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

		if ievent(event).Mask&unix.IN_ISDIR != 0 {
			continue // Skip directories
		}

		if !w.validFile(event.Path()) {
			continue // Skip invalid files
		}

		updateEvent := convertEvent(event.Event())
		if w.suspendEvent > 0 && updateEvent == w.suspendEvent {
			w.suspendEvent = 0
			log.Debugf("FileWatcher: Skipping suspended event %s for path: %q", updateEvent, event.Path())
			continue // Skip the suspended event
		}

		// Get any events registered for the specific file, and append the specified event
		var eventList notifyEvents
		if val, ok := w.batcher.Load(event.Path()); ok {
			eventList = val.(notifyEvents)
		}

		eventList = append(eventList, event)

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
			// Concatenate all known events, and dispatch them to be handled one by one
			for _, event := range w.concatenateEvents(val.(notifyEvents)) {
				w.sendUpdate(event)
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

func (w *FileWatcher) sendUpdate(update *FileUpdate) {
	log.Debugf("FileWatcher: Sending update: %s -> %q", update.Event, update.Path)
	w.updates <- update
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
func (w *FileWatcher) Suspend(updateEvent FileEvent) {
	w.suspendEvent = updateEvent
}

// validSuffix is used to filter out all unsupported
// files based on if their extension is unknown or
// if their path contains an excluded directory
func (w *FileWatcher) validFile(path string) bool {
	parts := strings.Split(filepath.Clean(path), string(os.PathSeparator))
	ext := filepath.Ext(parts[len(parts)-1])
	for _, suffix := range w.opts.ValidExtensions {
		if ext == suffix {
			return true
		}
	}

	for i := 0; i < len(parts)-1; i++ {
		for _, exclude := range w.opts.ExcludeDirs {
			if parts[i] == exclude {
				return false
			}
		}
	}

	return false
}

func convertEvent(event notify.Event) FileEvent {
	if updateEvent, ok := eventMap[event]; ok {
		return updateEvent
	}

	return FileEventNone
}

func convertUpdate(event notify.EventInfo) *FileUpdate {
	fileEvent := convertEvent(event.Event())
	if fileEvent == FileEventNone {
		// This should never happen
		panic(fmt.Sprintf("invalid event for update conversion: %q", event.Event().String()))
	}

	return &FileUpdate{
		Event: fileEvent,
		Path:  event.Path(),
	}
}

// moveCache caches an event during a move operation
// and dispatches a FileUpdate if it's not cancelled
type moveCache struct {
	watcher *FileWatcher
	event   notify.EventInfo
	timer   *time.Timer
}

func (w *FileWatcher) newMoveCache(event notify.EventInfo) *moveCache {
	m := &moveCache{
		watcher: w,
		event:   event,
	}

	// moveCaches wait one second to be cancelled before firing
	m.timer = time.AfterFunc(time.Second, m.incomplete)
	return m
}

func (m *moveCache) cookie() uint32 {
	return ievent(m.event).Cookie
}

// If the moveCache isn't cancelled, the move is considered incomplete and this
// method is fired. A complete move consists out of a "from" event and a "to" event,
// if only one is received, the file is moved in/out of a watched directory, which
// is treated as a normal creation/deletion by this method.
func (m *moveCache) incomplete() {
	var event FileEvent

	switch m.event.Event() {
	case notify.InMovedFrom:
		event = FileEventDelete
	case notify.InMovedTo:
		event = FileEventModify
	default:
		// This should never happen
		panic(fmt.Sprintf("moveCache: unrecognized event: %v", m.event.Event()))
	}

	log.Tracef("moveCache: Timer expired for %d, dispatching...", m.cookie())
	m.watcher.sendUpdate(&FileUpdate{event, m.event.Path()})

	// Delete the cache after the timer has fired
	delete(moveCaches, m.cookie())
}

func (m *moveCache) cancel() {
	m.timer.Stop()
	delete(moveCaches, m.cookie())
	log.Tracef("moveCache: Dispatching cancelled for %d", m.cookie())
}

// moveCaches keeps track of active moves by cookie
var moveCaches = make(map[uint32]*moveCache)

// move processes InMovedFrom and InMovedTo events in any order
// and dispatches FileUpdates when a move is detected
func (w *FileWatcher) move(event notify.EventInfo) (moveUpdate *FileUpdate) {
	cookie := ievent(event).Cookie
	cache, ok := moveCaches[cookie]
	if !ok {
		// The cookie is not cached, create a new cache object for it
		moveCaches[cookie] = w.newMoveCache(event)
		return
	}

	sourcePath, destPath := cache.event.Path(), event.Path()
	switch event.Event() {
	case notify.InMovedFrom:
		sourcePath, destPath = destPath, sourcePath
		fallthrough
	case notify.InMovedTo:
		cache.cancel()                                    // Cancel dispatching the cache's incomplete move
		moveUpdate = &FileUpdate{FileEventMove, destPath} // Register an internal, complete move instead
		log.Tracef("FileWatcher: Detected move: %q -> %q", sourcePath, destPath)
	}

	return
}

// concatenateEvents takes in a slice of events and concatenates
// all events possible based on combinedEvents. It also manages
// file moving and conversion from notifyEvents to FileUpdates
func (w *FileWatcher) concatenateEvents(events notifyEvents) FileUpdates {
	for _, combinedEvent := range combinedEvents {
		// Test if the prefix of the given events matches combinedEvent.input
		if event, ok := combinedEvent.match(events); ok {
			// If so, replace combinedEvent.input prefix in events with combinedEvent.output and recurse
			concatenated := events[len(combinedEvent.input):]
			if event != nil { // Prepend the concatenation result event if any
				concatenated = append(notifyEvents{event}, concatenated...)
			}

			log.Tracef("FileWatcher: Concatenated events: %v -> %v", events, concatenated)
			return w.concatenateEvents(concatenated)
		}
	}

	// Convert the events to updates
	updates := make(FileUpdates, 0, len(events))
	for _, event := range events {
		switch event.Event() {
		case notify.InMovedFrom, notify.InMovedTo:
			// Send move-related events to w.move
			if update := w.move(event); update != nil {
				// Add the update to the list if we get something back
				updates = append(updates, update)
			}
		default:
			updates = append(updates, convertUpdate(event))
		}
	}

	return updates
}

func ievent(event notify.EventInfo) *unix.InotifyEvent {
	return event.Sys().(*unix.InotifyEvent)
}
