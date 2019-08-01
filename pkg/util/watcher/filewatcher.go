package watcher

import (
	"bytes"
	"fmt"
	"github.com/weaveworks/ignite/pkg/util"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tywkeene/go-fsevents"
	"github.com/weaveworks/ignite/pkg/util/sync"
)

const updateBuffer = 4096 // How many updates we can buffer before watching is interrupted
var watchMask = fsevents.DirCreatedEvent | fsevents.DirRemovedEvent | fsevents.CloseWrite

var eventMap = map[uint32]Event{
	fsevents.FileCreatedEvent: EventCreate,
	fsevents.FileRemovedEvent: EventDelete,
	fsevents.CloseWrite:       EventModify,
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

type FileUpdateStream chan *FileUpdate

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

// FileWatcher recursively monitors changes in files in the given directory
// and sends out events based on their state changes. Only files conforming
// to validSuffix are monitored. The FileWatcher can be suspended for a single
// event at a time to eliminate updates by WatchStorage causing a loop.
type FileWatcher struct {
	dir          string
	updates      FileUpdateStream
	suspendEvent Event
	watcher      *fsevents.Watcher
	monitor      *sync.Monitor
	dispatcher   *sync.Monitor
	opts         Options
	// the batcher is used for properly sending many concurrent inotify events
	// as a group, after a specified timeout. This fixes the issue of one single
	// file operation being registered as many different inotify events
	batcher *sync.BatchWriter
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
		updates: make(FileUpdateStream, updateBuffer),
		batcher: sync.NewBatchWriter(opts.BatchTimeout),
		opts:    opts,
	}

	if w.watcher, err = fsevents.NewWatcher(dir, watchMask); err != nil {
		return
	}

	if err = w.start(w.dir, &files); err != nil {
		if err2 := w.watcher.StopAll(); err2 != nil {
			err = fmt.Errorf("%v, error stopping descriptors: %v", err, err2)
		}
	} else {
		go w.watcher.Watch()
		w.monitor = sync.RunMonitor(w.monitorFunc)
		w.dispatcher = sync.RunMonitor(w.dispatchFunc)
	}

	return
}

func (w *FileWatcher) addWatch(path string) (err error) {
	log.Tracef("FileWatcher: Adding watch for %q", path)
	wd, err := w.watcher.AddDescriptor(path, watchMask)
	if err == nil {
		err = wd.Start()
	}

	return
}

func (w *FileWatcher) hasWatch(path string) (b bool) {
	if b = w.watcher.DescriptorExists(path); b {
		log.Tracef("FileWatcher: Watch found for %q", path)
	} else {
		log.Tracef("FileWatcher: No watch found for %q", path)
	}

	return
}

func (w *FileWatcher) removeWatch(path string) error {
	log.Tracef("FileWatcher: Removing watch for %q", path)
	return w.watcher.RemoveDescriptor(path)
}

// start discovers all subdirectories and adds paths to
// notify before starting the monitoring goroutine
func (w *FileWatcher) start(dir string, files *[]string) error {
	//descriptors := w.watcher.ListDescriptors()
	//validDescriptors := make(map[string]bool, len(descriptors))
	//for _, path := range descriptors {
	//	validDescriptors[path] = false
	//}

	if err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip the root directory
			if path == w.dir {
				return nil
			}

			if info.IsDir() {
				for _, dir := range w.opts.ExcludeDirs {
					if info.Name() == dir {
						return filepath.SkipDir // Skip excluded directories
					}
				}

				if !w.watcher.DescriptorExists(path) {
					return w.addWatch(path)
				}

				//validDescriptors[path] = true
				//return nil
			}

			if files != nil {
				// Only include files with a valid suffix
				if w.validSuffix(info.Name()) {
					*files = append(*files, path)
				}
			}

			return nil
		}); err != nil {
		return err
	}

	//for path, valid := range validDescriptors {
	//	if !valid {
	//		if err := w.removeWatch(path); err != nil {
	//			return err
	//		}
	//	}
	//}

	return nil
}

func (w *FileWatcher) stop(dir string) (err error) {
	// TODO: On removal this might throw fsevents.ErrDescForEventNotFound, fix it

	if !w.watcher.DescriptorExists(dir) {
		return // Bail if this directory's descriptor is already gone
	}

	if !util.Subpath(w.dir, dir) {
		return fmt.Errorf("tried to remove watch from %q", w.dir)
	}

	parentDir := filepath.Dir(dir)
	if !util.DirExists(parentDir) {
		// If the parent directory has already been deleted, call this method
		// recursively until hitting the lowest directory that was deleted
		// and perform the descriptor lookup and removal from there
		return w.stop(parentDir)
	}

	for _, path := range w.watcher.ListDescriptors() {
		if util.Subpath(dir, path) {
			if err = w.removeWatch(path); err != nil {
				return
			}
		}
	}

	fmt.Println("Number of watches:", len(w.watcher.ListDescriptors()))

	return
}

func (w *FileWatcher) monitorFunc() {
	log.Debug("FileWatcher: Monitoring thread started")
	defer log.Debug("FileWatcher: Monitoring thread stopped")
	defer close(w.updates) // Close the update stream after the FileWatcher has stopped

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			log.Infoln("Event:", visualize(event))

			//fmt.Println("Watches before:", len(w.watcher.ListDescriptors()))

			//if fsevents.CheckMask(fsevents.IsDir | fsevents.RootDelete, event.RawEvent.Mask) {
			//	fmt.Println("ROOT DIR delete:", event.Path)
			//}

			// If the event targets a directory, handle the watchers for it
			if w.handleDirEvent(event) {
				//fmt.Println("Watches after:", len(w.watcher.ListDescriptors()))
				continue // Skip file processing for the event
			}

			updateEvent := convertEvent(event)
			if w.suspendEvent > 0 && updateEvent == w.suspendEvent {
				w.suspendEvent = 0
				log.Debugf("FileWatcher: Skipping suspended event %s for path: %q", updateEvent, event.Path)
				continue // Skip the suspended event
			}

			// Suppress successive duplicate events registered in suppressDuplicates
			//if suppressEvent(event.Path, updateEvent) {
			//	log.Debugf("FileWatcher: Skipping suppressed event %s for path: %q", updateEvent, event.Path)
			//	continue // Skip the suppressed event
			//}

			// Get any events registered for the specific file, and append the specified event
			var eventList Events
			if val, ok := w.batcher.Load(event.Path); ok {
				eventList = val.(Events)
			}

			eventList = append(eventList, updateEvent)

			// Register the event in the map, and dispatch all the events at once after the timeout
			w.batcher.Store(event.Path, eventList)
			log.Debugf("FileWatcher: Registered inotify events %v for path %q", eventList, event.Path)
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}

			/*
				When fsevents reads a single event using C.read, it's source watch
				is parsed to a descriptor using w.GetDescriptorByWatch after the read
				has completed. If the descriptor gets deleted by w.stop during this time
				(based on the directory getting deleted), w.GetDescriptorByWatch fails,
				fsevents.ErrDescForEventNotFound is returned and no event is fired.
				This error is safe to ignore, as we don't want any events referencing
				deleted directories and it doesn't cause any event misses as the watch
				has already been deleted by other means.
			*/
			if err != fsevents.ErrDescForEventNotFound {
				log.Errorf("FileWatcher: Error watching: %v", err)
			}
		}
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

func (w *FileWatcher) handleDirEvent(event *fsevents.FsEvent) bool {
	//if err := w.start(nil); err != nil {
	//	log.Errorf("FileWatcher: Watch directory updating failed: %v", err)
	//}

	if event.IsDirCreated() {
		if err := w.start(event.Path, nil); err != nil {
			log.Errorf("FileWatcher: Failed to add watch %q: %v", event.Path, err)
		}
		//} else if fsevents.CheckMask(fsevents.RootDelete, event.RawEvent.Mask) {
	} else if event.IsDirRemoved() {
		if err := w.stop(event.Path); err != nil {
			log.Errorf("FileWatcher: Failed remove watch %q: %v", event.Path, err)
		}
	} else {
		return false
	}

	return true
}

// GetFileUpdateStream gets the channel with FileUpdates
func (w *FileWatcher) GetFileUpdateStream() FileUpdateStream {
	return w.updates
}

// Close closes active underlying resources
func (w *FileWatcher) Close() {
	//w.watcher.StopAll()
	w.batcher.Close()
	close(w.watcher.Events)
	close(w.watcher.Errors)
	//close(w.events) // Close the event stream
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

func convertEvent(event *fsevents.FsEvent) Event {
	for mask, e := range eventMap {
		if fsevents.CheckMask(mask, event.RawEvent.Mask) {
			return e
		}
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

var visualMap = map[uint32]string{
	fsevents.Accessed:   "IN_ACCESS",
	fsevents.Modified:   "IN_MODIFY",
	fsevents.AttrChange: "IN_ATTRIB",
	fsevents.CloseWrite: "IN_CLOSE_WRITE",
	fsevents.CloseRead:  "IN_CLOSE_NOWRITE",
	fsevents.Open:       "IN_OPEN",
	fsevents.MovedFrom:  "IN_MOVED_FROM",
	fsevents.MovedTo:    "IN_MOVED_TO",
	fsevents.Move:       "IN_MOVE",
	fsevents.Create:     "IN_CREATE",
	fsevents.Delete:     "IN_DELETE",
	fsevents.RootDelete: "IN_DELETE_SELF",
	fsevents.RootMove:   "IN_MOVE_SELF",
	fsevents.IsDir:      "IN_ISDIR",
}

var visualOrder = []uint32{
	fsevents.Accessed,
	fsevents.Modified,
	fsevents.AttrChange,
	fsevents.CloseWrite,
	fsevents.CloseRead,
	fsevents.Open,
	fsevents.MovedFrom,
	fsevents.MovedTo,
	fsevents.Move,
	fsevents.Create,
	fsevents.Delete,
	fsevents.RootDelete,
	fsevents.RootMove,
	fsevents.IsDir,
}

// TODO: This is a helper for visualizing go-fsevents events
func visualize(event *fsevents.FsEvent) string {
	var sb strings.Builder

	for _, bit := range visualOrder {
		if fsevents.CheckMask(bit, event.RawEvent.Mask) {
			if sb.Len() > 0 {
				sb.WriteString(" | ")
			}

			sb.WriteString(visualMap[bit])
		}
	}

	return sb.String()
}
