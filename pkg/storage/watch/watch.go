package watch

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/storage"
	"os"
	"path/filepath"
	"syscall"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/storage/watch/update"
)

const eventBuffer = 4096 // How many events we can buffer before watching is interrupted
var excludeDirs = []string{".git"}

var eventMap = map[fsnotify.Op]update.Event{
	fsnotify.Create: update.EventCreate,
	fsnotify.Remove: update.EventDelete,
	fsnotify.Write:  update.EventModify,
}

type EventStream chan *update.FileUpdate

type watcher struct {
	dir          string
	events       EventStream
	watcher      *fsnotify.Watcher
	suspendEvent update.Event
	monitor      *Monitor
}

// newWatcher returns a list of files in the watched directory in
// addition to the generated watcher, it can be used to populate
// MappedRawStorage fileMappings
func newWatcher(dir string) (w *watcher, files []string, err error) {
	w = &watcher{
		dir:    dir,
		events: make(EventStream, eventBuffer),
	}

	w.watcher, err = fsnotify.NewWatcher()
	if err == nil {
		files, err = w.start()
	}

	if err != nil {
		if closeErr := w.watcher.Close(); closeErr != nil {
			err = fmt.Errorf("%v, error closing: %v", err, closeErr)
		}
	}

	return
}

// start discovers all subdirectories and adds paths to
// fsnotify before starting the monitoring goroutine
func (w *watcher) start() (files []string, err error) {
	if err = filepath.Walk(w.dir,
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

				return w.watcher.Add(path)
			}

			// Only include files with a valid suffix
			if validSuffix(info.Name()) {
				files = append(files, path)
			}

			return nil
		}); err == nil {
		w.monitor = RunMonitor(w.monitorFunc)
	}

	return
}

func (w *watcher) monitorFunc() {
	defer log.Debug("Watcher: monitoring thread stopped")
	defer close(w.events) // Close the event stream after the watcher has stopped

	for {
		select {
		case e, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			event := convertEvent(e)
			if event == 0 {
				log.Debug("Watcher: skipping unregistered event")
				continue // Skip all unregistered events
			}

			if event == w.suspendEvent {
				log.Debugf("Watcher: skipping suspended event: %s", event)
				continue // Skip suspended events
			}

			log.Debugf("Watcher: event: %s", event)

			if event == update.EventDelete {
				if err := w.watcher.Remove(e.Name); err != nil {
					// Watcher removal will fail with textual output if it doesn't
					// exist, if an actual error occurs, it's outputted as syscall.Errno
					if errno, ok := err.(syscall.Errno); ok {
						log.Errorf("Watcher: failed to remove %q: %v", e.Name, errno)
					}
				}
			} else if event == update.EventCreate {
				if fi, err := os.Stat(e.Name); err != nil {
					log.Errorf("Watcher: failed to stat %q: %v", e.Name, err)
				} else {
					if fi.IsDir() {
						if err := w.watcher.Add(e.Name); err != nil {
							log.Errorf("Watcher: failed to add %q: %v", e.Name, err)
						}

						continue
					}
				}
			}

			// Only fire events for files with valid extensions: .yaml, .yml, .json
			if validSuffix(e.Name) {
				log.Debugf("Watcher: sending update: %s -> %q", event, e.Name)
				w.events <- &update.FileUpdate{
					Event: event,
					Path:  e.Name,
				}
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}

			log.Errorf("Watcher: error: %v", err)
		}
	}
}

func (w *watcher) close() {
	_ = w.watcher.Close() // This returns only nil errors
	w.monitor.Wait()
}

func (w *watcher) suspend(event update.Event) {
	w.suspendEvent = event
}

func (w *watcher) resume() {
	w.suspendEvent = 0
}

func convertEvent(event fsnotify.Event) update.Event {
	for fsEvent, updateEvent := range eventMap {
		if event.Op&fsEvent == fsEvent {
			return updateEvent
		}
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
