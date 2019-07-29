package watch

import (
	"os"
	"path/filepath"

	"github.com/rjeczalik/notify"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/storage"
	"github.com/weaveworks/ignite/pkg/storage/watch/update"
)

const eventBuffer = 4096 // How many events and updates we can buffer before watching is interrupted
var excludeDirs = []string{".git"}
var listenEvents = []notify.Event{notify.InCreate, notify.InDelete, notify.InDeleteSelf, notify.InCloseWrite}

var eventMap = map[notify.Event]update.Event{
	notify.InCreate:     update.EventCreate,
	notify.InDelete:     update.EventDelete,
	notify.InCloseWrite: update.EventModify,
}

type eventStream chan notify.EventInfo
type UpdateStream chan *update.FileUpdate
type watches []string

type watcher struct {
	dir          string
	events       eventStream
	updates      UpdateStream
	watches      watches
	suspendEvent update.Event
	monitor      *Monitor
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
	}

	if err = w.start(&files); err != nil {
		notify.Stop(w.events)
	} else {
		w.monitor = RunMonitor(w.monitorFunc)
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

		log.Debugf("Watcher: event: %s", event)

		switch event.Event() {
		case notify.InCreate:
			if fi, err := os.Stat(event.Path()); err != nil {
				log.Errorf("Watcher: failed to stat %q: %v", event.Path(), err)
				continue
			} else {
				if fi.IsDir() {
					if err := w.addWatch(event.Path()); err != nil {
						log.Errorf("Watcher: failed to add %q: %v", event.Path(), err)
					}

					continue
				}
			}

			fallthrough
		case notify.InDelete, notify.InCloseWrite:
			if event.Event() == notify.InDelete && w.hasWatch(event.Path()) {
				w.clear()
				if err := w.start(nil); err != nil {
					log.Errorf("Watcher: Failed to re-initialize watches for %q", w.dir)
				}

				continue
			}

			if validSuffix(event.Path()) {
				if updateEvent > 0 {
					log.Debugf("Watcher: sending update: %s -> %q", updateEvent, event.Path())
					w.updates <- &update.FileUpdate{
						Event: updateEvent,
						Path:  event.Path(),
					}
				}
			}
		}
	}
}

func (w *watcher) close() {
	notify.Stop(w.events)
	close(w.events) // Close the event stream
	w.monitor.Wait()
}

func (w *watcher) suspend(updateEvent update.Event) {
	w.suspendEvent = updateEvent
}

func convertEvent(event notify.Event) update.Event {
	if updateEvent, ok := eventMap[event]; ok {
		return updateEvent
	}

	return 0
}

func validSuffix(path string) bool {
	for suffix := range storage.Formats {
		if filepath.Ext(path) == suffix {
			return true
		}
	}

	return false
}
