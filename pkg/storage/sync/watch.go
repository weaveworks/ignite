package sync

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/storage"
	"path/filepath"
)

type watcher struct {
	storage    storage.Storage
	updateFunc updateFunc
	watcher    *fsnotify.Watcher
}

type updateFunc func(*fileUpdate) error

func newWatcher(storage storage.Storage, updateFunc updateFunc) (w *watcher, err error) {
	w = &watcher{
		storage:    storage,
		updateFunc: updateFunc,
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

func (w *watcher) start() error {
	if err := filepath.Walk(w.storage.RawStorage().Dir(),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return w.watcher.Add(path)
			}

			return nil
		}); err != nil {
		return err
	}

	go func() {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return // TODO: Close the watcher?
			}

			log.Debugf("watcher: event: %v", event)

			fu := &fileUpdate{
				storage: w.storage,
				path:    event.Name,
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				fu.event = eventModify
			} else if event.Op&fsnotify.Create == fsnotify.Create {
				if fi, err := os.Stat(event.Name); err != nil {
					log.Errorf("watcher: failed to stat %q: %v", event.Name, err)
				} else {
					if fi.IsDir() {
						if err := w.watcher.Add(event.Name); err != nil {
							log.Errorf("watcher: failed to add %q: %v", event.Name, err)
						}
					} else {
						fu.event = eventCreate
					}
				}
			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				if fi, err := os.Stat(event.Name); err != nil {
					log.Errorf("watcher: failed to stat %q:", event.Name, err)
				} else {
					if fi.IsDir() {
						if err := w.watcher.Remove(event.Name); err != nil {
							log.Errorf("watcher: failed to remove %q:", event.Name, err)
						}
					} else {
						fu.event = eventDelete
					}
				}
			}

			if fu.event > 0 {
				if err := w.updateFunc(fu); err != nil {
					log.Errorf("watcher: updating %q failed: %v", event.Name, err)
				}
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return // TODO: Close the watcher?
			}

			log.Errorf("watcher: error: %v", err)
		}
	}()
}
