package watch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/storage/watch/update"
)

const eventBuffer = 4096 // How many events we can buffer before watching is interrupted

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
}

func newWatcher(dir string) (w *watcher, err error) {
	w = &watcher{
		dir:    dir,
		events: make(EventStream, eventBuffer),
	}

	w.watcher, err = fsnotify.NewWatcher()
	if err == nil {
		err = w.start()
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
func (w *watcher) start() (err error) {
	if err = filepath.Walk(w.dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return w.watcher.Add(path)
			}

			return nil
		}); err == nil {
		go w.monitor()
	}

	return
}

func (w *watcher) monitor() {
	defer func() {
		if err := w.watcher.Close(); err != nil {
			log.Errorf("Failed to close watcher: %v", err)
		}
	}()

	for {
		select {
		case e, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			event := convertEvent(e)
			if event == 0 {
				log.Debugf("Watcher: skipping unregistered event: %v", event)
				continue // Skip all unregistered events
			}

			if event == w.suspendEvent {
				log.Debugf("Watcher: skipping suspended event: %v", event)
				continue // Skip suspended events
			}

			log.Debugf("Watcher: event: %v", event)

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
	filename := filepath.Base(path)

	for _, suffix := range []string{".json", ".yaml", ".yml"} {
		if strings.HasSuffix(filename, suffix) {
			return true
		}
	}

	return false
}
