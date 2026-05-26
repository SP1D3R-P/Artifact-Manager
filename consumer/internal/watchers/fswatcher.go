/*
* File system watcher implementation using fsnotify. It provides an interface to watch
* for file system events and execute callbacks when those events occur.
 */
package watchers

import (
	"context"
	"log"

	"github.com/fsnotify/fsnotify"
)

type FSCallbackFn interface {
	Call(arg interface{}) error
}

type FSWatcher interface {
	AddPaths(paths ...string) error
	AddCallbacks(event fsnotify.Op, callback ...FSCallbackFn) error
	Watch(ctx context.Context)
	Close()
}

type watcher struct {
	fsWatcher *fsnotify.Watcher
	callbacks map[fsnotify.Op][]FSCallbackFn
}

func NewFSWatcher() (FSWatcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &watcher{
		fsWatcher: fsWatcher,
		callbacks: make(map[fsnotify.Op][]FSCallbackFn),
	}, nil
}

func (w *watcher) AddPaths(paths ...string) error {
	for _, path := range paths {
		if err := w.fsWatcher.Add(path); err != nil {
			return err
		}
	}
	return nil
}

func (w *watcher) Watch(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return

		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}

			for _, callback := range w.callbacks[event.Op] {
				if err := callback.Call(event.Name); err != nil {
					log.Printf("Callback error for event %s: %v", event.Op, err)
				}
			}

		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}

}

func (w *watcher) AddCallbacks(event fsnotify.Op, callback ...FSCallbackFn) error {
	w.callbacks[event] = append(w.callbacks[event], callback...)
	return nil
}

func (w *watcher) Close() {
	w.fsWatcher.Close()
}
