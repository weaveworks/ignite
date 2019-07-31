# go-fsevents
[![GoDoc](https://godoc.org/github.com/tywkeene/go-fsevents?status.svg)](https://godoc.org/github.com/tywkeene/go-fsevents)
[![Build Status](https://travis-ci.org/tywkeene/go-fsevents.svg?branch=master)](https://travis-ci.org/tywkeene/go-fsevents)
[![codecov.io Code Coverage](https://img.shields.io/codecov/c/github/tywkeene/go-fsevents.svg?maxAge=2592000)](https://codecov.io/github/tywkeene/go-fsevents?branch=master)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)  
[![Go Report Card](https://goreportcard.com/badge/github.com/tywkeene/go-fsevents)](https://goreportcard.com/report/github.com/tywkeene/go-fsevents)


Recursive filesystem event watcher using inotify in golang

go-fsevents provides functions necessary for monitoring filesystem events on linux systems using the inotify interface.

Unlike other inotify packages, go-fsevents provides a recursive watcher, allowing the monitoring of directory trees easily.

## UNSTABLE

The package is currently unstable, and as such should not be used in any production environment.

Many changes, additions and breaking refactors will take place between now and the stable 1.0.0 release.

You have been warned.

## Features

- Single directory event monitoring
- Recursive directory tree event monitoring
- EventHandle interface to allow for clean and concise handling of events
- Access to the underlying raw inotify event through the [unix](https://godoc.org/golang.org/x/sys/unix) package
- Predefined event translations. No need to fuss with raw inotify flags.
- Concurrency safe


## Quickstart

To use go-fsevents we simply need a Watcher and a mask that will describe the events we want to watch, In this example we'll watch for file and directory creation events.

```go
mask := fsevents.FileCreated | fsevents.DirectoryCreated
w, err := fsevents.NewWatcher("foo/bar", mask)
```

Then we can start the watcher and the go-routine that will write events to a channel

```go
w.StartAll()
go w.Watch()
```

Now we can read from the Watcher's channels

```go
for {
    select {
           case event := <-w.Events:
               // Here we can also add, start and remove descriptors in response to events
               if event.IsDirCreated() == true {
                   descriptor, err := w.AddDescriptor(event.Path, w.Mask)
                   if err != nil{
                       log.Println("Error adding descriptor:", err)
                       break
                   }
                   descriptor.Start()
               }
           
           case err := <-w.Errors:
               log.Println("Error:", err)
               break
    }
}
```

Now any directory under `foo/bar` will be watched, and events in these directories will be written to the w.Events channel

When we're done with a descriptor, we can stop it

```go
descriptor := w.GetDescriptorByPath("foo/bar/baz")
descriptor.Stop()
```

Or if we want to stop all of the watches

```go
w.StopAll()
```

## Handles

The `EventHandlers` interface provides a clean and easy way to automatically and pragmatically handle events:

The interface is as such

```go
// EventHandler allows for the Watcher to apply pre-registered functions in response to an event.
type EventHandler interface {
	// The Handle method is called by WatchAndHandle in response to a given event
	Handle(w *Watcher, event *FsEvent) error
	// The Check method is called to match Events with the correct EventHandle in the Watcher
	// Check must return true if the event described by the in event matches the argument
	Check(event *FsEvent) bool
	// The GetMask method returns the uint32 inotify mask this EventHandle handles
	GetMask() uint32
}
```

We create a struct that satisfies this interface, and register it with the mask it will handle in a `Watcher` instance

```go
// The DirectoryCreatedHandle implements the EventHandler interface
type DirectoryCreatedHandle struct {
	Mask uint32
}

// The Handle has access Watcher and Event this Handle was called from
// Returning an error from an EventHandler causes the Watcher to write the error to the Error channel
func (h *DirectoryCreatedHandle) Handle(w *fsevents.Watcher, event *fsevents.FsEvent) error {
	log.Println("Directory created:", event.Path)

	// The Watcher can be used inside event handles to add/remove/modify Watches
	// In this case, we add a descriptor for the created directory and start a watch for it
	d, err := w.AddDescriptor(event.Path, h.GetMask())
	if err != nil {
		return err
	}

	d.Start()
	log.Printf("Started watch on %q", event.Path)
	return nil
}

// GetMask returns the inotify event mask this EventHandler handles
func (h *DirectoryCreatedHandle) GetMask() uint32 {
	return h.Mask
}

// The most basic usage is to use the FsEvent methods to check the event mask against the handler
// Other logic may be used to determine if the handle should be executed for the given event.
func (h *DirectoryCreatedHandle) Check(event *fsevents.FsEvent) bool {
	return event.IsDirCreated()
}

```

Then we can register the handler with the Watcher instance

```go
	w.RegisterEventHandler(&DirectoryCreatedHandle{Mask: fsevents.DirCreatedEvent})
```

And now start the watches and go-routine that will execute the registered handler in response to events

```go
	w.StartAll()

	// WatchAndHandle will search for the correct handle in response to a given
	// event and apply it, writing the error it returns, if any, to w.Errors
	go w.WatchAndHandle()
```






## Examples

See the examples in [examples](https://github.com/tywkeene/go-fsevents/blob/master/examples) for quick and easy runnable examples of how go-fsevents can be used in your project

`handlers.go` describes how to use the `EventHandlers` interface to handle events automatically

`loop.go` describes how to read events from the `watcher.Events` channel
