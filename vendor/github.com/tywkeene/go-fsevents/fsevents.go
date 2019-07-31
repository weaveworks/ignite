// Package fsevents provides routines for monitoring filesystem events
// on a Linux system via the inotify subsystem recursively.
package fsevents

// #include <unistd.h>
import "C"

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// WatchDescriptor describes a path being watched
type WatchDescriptor struct {
	// The path of this descriptor
	Path string
	// This descriptor's inotify watch mask
	Mask uint32
	// This descriptor's inotify watch descriptor
	WatchDescriptor int
	// Is this watcher currently running?
	Running bool
	// InotifyDescriptor of the Watcher this WatchDescriptor belongs to
	InotifyDescriptor *int
}

// FsEvent is an inotify event along with the ID and timestamp of the event
type FsEvent struct {
	// The name of the event's file
	Name string
	// The full path of the event
	Path string
	// The raw inotify event
	RawEvent *unix.InotifyEvent
	// The actual inotify watch descriptor related to this event
	Descriptor *WatchDescriptor
	// The serial ID of this event. ID is incremented in ReadSingleEvent upon successful event read
	ID uint32
	// Timestamp of the time the event occurred in UTC
	Timestamp time.Time
}

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

// Watcher is the top-level object that handles EventHandlers and Watchers.
type Watcher struct {
	sync.Mutex
	// List of EventHandles that have been registered with this Watcher
	eventHandlers []EventHandler
	// The root path of this watcher
	RootPath string
	// The main inotify descriptor
	InotifyDescriptor int
	// Watch descriptors in this watch key: watch path -> value: WatchDescriptor
	Descriptors map[string]*WatchDescriptor
	// How many events have been read by this watcher from the inotify descriptor
	// This counter is incremented in ReadSingleEvent
	EventCount uint32
	// The event channel we send all events on
	Events chan *FsEvent
	// How we report errors
	Errors chan error
}

var (
	// All the errors returned by fsevents
	// Should probably provide a more situationally descriptive message along with it

	//Top-level Watcher errors
	ErrWatchNotCreated      = errors.New("watcher could not be created")
	ErrNoRunningDescriptors = errors.New("watcher has no running descriptors")
	ErrNoEventHandles       = errors.New("watcher has no registered event handles")

	//Descriptor errors
	ErrDescNotCreated       = errors.New("descriptor could not be created")
	ErrDescNotStart         = errors.New("descriptor could not be started")
	ErrDescRunning          = errors.New("descriptor already running")
	ErrDescNotStopped       = errors.New("descriptor could not be stopped")
	ErrDescAlreadyExists    = errors.New("descriptor for that directory already exists")
	ErrDescNotRunning       = errors.New("descriptor not running")
	ErrDescForEventNotFound = errors.New("descriptor for event not found")
	ErrDescNotFound         = errors.New("descriptor not found")

	// Event handle errors
	ErrNoSuchHandle = errors.New("handle not found")
	ErrHandleError  = errors.New("handle returned error")
	ErrHandleExists = errors.New("a handle for this mask already exists")

	//Inotify interface errors
	ErrIncompleteRead = errors.New("incomplete event read")
	ErrReadError      = errors.New("error reading an event")
)

var (
	// Default inotify flags
	Accessed   uint32 = unix.IN_ACCESS
	Modified   uint32 = unix.IN_MODIFY
	AttrChange uint32 = unix.IN_ATTRIB
	CloseWrite uint32 = unix.IN_CLOSE_WRITE
	CloseRead  uint32 = unix.IN_CLOSE_NOWRITE
	Open       uint32 = unix.IN_OPEN
	MovedFrom  uint32 = unix.IN_MOVED_FROM
	MovedTo    uint32 = unix.IN_MOVED_TO
	Move       uint32 = unix.IN_MOVE
	Create     uint32 = unix.IN_CREATE
	Delete     uint32 = unix.IN_DELETE
	RootDelete uint32 = unix.IN_DELETE_SELF
	RootMove   uint32 = unix.IN_MOVE_SELF
	IsDir      uint32 = unix.IN_ISDIR

	AllEvents = (Accessed | Modified | AttrChange | CloseWrite | CloseRead | Open | MovedFrom |
		MovedTo | MovedTo | Create | Delete | RootDelete | RootMove | IsDir)

	// Custom event flags

	// Directory events

	// A quick breakdown, same goes for the file events, except
	// those pertain to files, not directories. There is a difference.

	// The directory is not in the watch directory anymore
	// whether it was moved or deleted, it's *poof* gone
	DirRemovedEvent = MovedFrom | Delete | IsDir

	// Whether it was moved or copied into the watch directory,
	// or created with mkdir, there is a new directory
	DirCreatedEvent = MovedTo | Create | IsDir

	// A directory was closed with write permissions, modified, or its
	// attributes changed in some way
	DirChangedEvent = CloseWrite | Modified | AttrChange | IsDir

	// File events
	FileRemovedEvent = MovedFrom | Delete
	FileCreatedEvent = MovedTo | Create
	FileChangedEvent = CloseWrite | Modified | AttrChange

	// Root watch directory events
	RootEvent = RootDelete | RootMove
)

// CheckMask returns true if flag 'check' is found in bitmask 'mask'
func CheckMask(check uint32, mask uint32) bool {
	return (mask & check) != 0
}

// IsDirEvent Returns true if the event is a directory event
func (e *FsEvent) IsDirEvent() bool {
	return CheckMask(IsDir, e.RawEvent.Mask)
}

// Root events.

// IsRootDeletion returns true if the event contains the inotify flag IN_DELETE_SELF
// and the rootPath argument matches the path in the FsEvent structure.
// This means the root watch directory has been deleted,
// and there will be no more events read from the descriptor
// since it doesn't exist anymore. You should probably handle this
// gracefully and always check for this event before doing anything else
// Also be sure to add the RootDelete flag to your watched events when
// initializing fsevents
func (e *FsEvent) IsRootDeletion(rootPath string) bool {
	return (CheckMask(RootDelete, e.RawEvent.Mask) == true) && (rootPath == e.Path)
}

// IsRootMoved returns true if the event contains the inotify flag IN_MOVE_SELF
// and the rootPath argument matches the path in the FsEvent structure.
// This means the root watch directory has been moved. This may not matter
// to you at all, and depends on how you deal with paths in your program.
// Still, you should check for this event before doing anything else.
func (e *FsEvent) IsRootMoved(rootPath string) bool {
	return (CheckMask(RootMove, e.RawEvent.Mask) == true) && (rootPath == e.Path)
}

// Custom directory events

// IsDirChanged returns true if the event describes a directory that
// was closed with write permissions, modified, or its attributes changed
func (e *FsEvent) IsDirChanged() bool {
	return ((CheckMask(CloseWrite, e.RawEvent.Mask) == true) && (e.IsDirEvent() == true)) ||
		((CheckMask(Modified, e.RawEvent.Mask) == true) && (e.IsDirEvent() == true)) ||
		((CheckMask(AttrChange, e.RawEvent.Mask) == true) && (e.IsDirEvent() == true))
}

// IsDirCreated returns true if the event describes a directory created
// within the root watch, or moved into the root watch directory
func (e *FsEvent) IsDirCreated() bool {
	return ((CheckMask(Create, e.RawEvent.Mask) == true) && (e.IsDirEvent() == true)) ||
		((CheckMask(MovedTo, e.RawEvent.Mask) == true) && (e.IsDirEvent() == true))
}

// IsDirRemoved returns true if the event describes a directory that was
//deleted or moved out of the root watch directory
func (e *FsEvent) IsDirRemoved() bool {
	return ((CheckMask(Delete, e.RawEvent.Mask) == true) && (e.IsDirEvent() == true)) ||
		((CheckMask(MovedFrom, e.RawEvent.Mask) == true) && (e.IsDirEvent() == true))
}

// Custom file events

// IsFileCreated returns true if the event describes a file that was moved into,
// or created within the root watch directory
func (e *FsEvent) IsFileCreated() bool {
	return (((CheckMask(Create, e.RawEvent.Mask) == true) && (e.IsDirEvent() == false)) ||
		((CheckMask(MovedTo, e.RawEvent.Mask) == true) && (e.IsDirEvent() == false)))
}

// IsFileRemoved returns true if the event describes a file
// was deleted or moved out of the root watch directory
func (e *FsEvent) IsFileRemoved() bool {
	return ((CheckMask(Delete, e.RawEvent.Mask) == true) && (e.IsDirEvent() == false) ||
		((CheckMask(MovedFrom, e.RawEvent.Mask) == true) && (e.IsDirEvent() == false)))
}

// IsFileChanged returns true if the event describes a file that was
// closed with write permissions, modified, or its attributes changed
func (e *FsEvent) IsFileChanged() bool {
	return ((CheckMask(CloseWrite, e.RawEvent.Mask) == true) && (e.IsDirEvent() == false)) ||
		((CheckMask(Modified, e.RawEvent.Mask) == true) && (e.IsDirEvent() == false)) ||
		((CheckMask(AttrChange, e.RawEvent.Mask) == true) && (e.IsDirEvent() == false))
}

func newWatchDescriptor(dirPath string, mask uint32, inotifyDescriptor int) *WatchDescriptor {
	return &WatchDescriptor{
		Path:              dirPath,
		WatchDescriptor:   -1,
		Mask:              mask,
		InotifyDescriptor: &inotifyDescriptor,
	}
}

// Start starts a WatchDescriptor inotify event watcher. If the descriptor is already running Start returns ErrDescRunning
func (d *WatchDescriptor) Start() error {
	var err error
	if d.Running == true {
		return ErrDescRunning
	}
	d.WatchDescriptor, err = unix.InotifyAddWatch(*d.InotifyDescriptor, d.Path, d.Mask)
	if d.WatchDescriptor == -1 || err != nil {
		d.Running = false
		return fmt.Errorf("%s: %s", ErrDescNotStart, err)
	}
	d.Running = true
	return nil
}

// Stop stops a running watch descriptor. If the descriptor is not running Stop returns ErrDescNotRunning
func (d *WatchDescriptor) Stop() error {
	if d.Running == false {
		return ErrDescNotRunning
	}
	_, err := unix.InotifyRmWatch(*d.InotifyDescriptor, uint32(d.WatchDescriptor))
	if err != nil {
		return fmt.Errorf("%s: %s", ErrDescNotStopped, err)
	}
	d.Running = false
	return nil
}

// DoesPathExist returns true if the path described by the descriptor exists
func (d *WatchDescriptor) DoesPathExist() bool {
	_, err := os.Lstat(d.Path)
	return os.IsExist(err)
}

// DescriptorExists returns true if a WatchDescriptor exists in Watcher w, false otherwise
func (w *Watcher) DescriptorExists(watchPath string) bool {
	w.Lock()
	defer w.Unlock()
	if _, exists := w.Descriptors[watchPath]; exists {
		return true
	}
	return false
}

// ListDescriptors returns a string array of all WatchDescriptors in w *Watcher
// Both started and stopped. To get a count of running watch descriptors, use GetRunningDescriptors
func (w *Watcher) ListDescriptors() []string {
	list := make([]string, 0)
	w.Lock()
	defer w.Unlock()
	for path, _ := range w.Descriptors {
		list = append(list, path)
	}
	return list
}

// RemoveDescriptor removes the WatchDescriptor with the path matching path
// from the watcher, and stops the inotify watcher
func (w *Watcher) RemoveDescriptor(path string) error {
	if w.DescriptorExists(path) == false {
		return ErrDescNotFound
	}
	w.Lock()
	defer w.Unlock()
	descriptor := w.Descriptors[path]
	if descriptor.DoesPathExist() == true {
		if err := descriptor.Stop(); err != nil {
			return err
		}
	}
	delete(w.Descriptors, path)
	return nil
}

// AddDescriptor adds a descriptor to Watcher w. The descriptor is not started.
func (w *Watcher) AddDescriptor(dirPath string, mask uint32) (*WatchDescriptor, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s: %s", ErrDescNotCreated, "directory does not exist")
	}
	if w.DescriptorExists(dirPath) == true {
		return nil, ErrDescAlreadyExists
	}

	descriptor := newWatchDescriptor(dirPath, mask, w.InotifyDescriptor)

	w.Lock()
	w.Descriptors[dirPath] = descriptor
	w.Unlock()

	return descriptor, nil
}

// RecursiveAdd adds the directory at rootPath, and all directories below it, using the flags provided in mask
func (w *Watcher) RecursiveAdd(rootPath string, mask uint32) error {
	dirStat, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return err
	}

	for _, child := range dirStat {
		if child.IsDir() == true {
			childPath := path.Clean(path.Join(rootPath, child.Name()))
			if err := w.RecursiveAdd(childPath, mask); err != nil {
				return fmt.Errorf("could not add recurisve-descriptor for path %q: %s", childPath, err.Error())
			}
			d, err := w.AddDescriptor(childPath, mask)
			if err != nil {
				return fmt.Errorf("could not add descriptor for path %q: %s", childPath, err.Error())
			}
			if err := d.Start(); err != nil {
				return fmt.Errorf("could not start watch for path %q: %s", childPath, err.Error())
			}
		}
	}
	return nil
}

// NewWatcher allocates a new watcher at path rootPath, and adds a descriptor with the mask provided
// This function initializes inotify, so it must be run first
func NewWatcher(rootPath string, mask uint32) (*Watcher, error) {
	fd, err := unix.InotifyInit()
	if fd == -1 || err != nil {
		return nil, fmt.Errorf("%s: %s", ErrWatchNotCreated, err)
	}

	w := &Watcher{
		eventHandlers:     make([]EventHandler, 0),
		RootPath:          path.Clean(rootPath),
		InotifyDescriptor: fd,
		Descriptors:       make(map[string]*WatchDescriptor),
		Events:            make(chan *FsEvent),
		Errors:            make(chan error),
	}

	_, err = w.AddDescriptor(w.RootPath, mask)

	return w, err
}

// GetRunningDescriptors returns the count of currently running or Start()'d descriptors for this watcher.
func (w *Watcher) GetRunningDescriptors() int32 {
	w.Lock()
	defer w.Unlock()
	var count int32
	for _, d := range w.Descriptors {
		if d.Running == true {
			count++
		}
	}
	return count
}

// StartAll starts all inotify watches described by this Watcher
func (w *Watcher) StartAll() error {
	w.Lock()
	defer w.Unlock()
	for _, d := range w.Descriptors {
		if err := d.Start(); err != nil {
			return err
		}
	}
	return nil
}

// StopAll stops all running watch descriptors. Does not remove descriptors from the watch
func (w *Watcher) StopAll() error {
	w.Lock()
	defer w.Unlock()
	for _, d := range w.Descriptors {
		if d.Running == true {
			if err := d.Stop(); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetDescriptorByWatch searches a Watcher instance for a watch descriptor.
// Searches by inotify watch descriptor
func (w *Watcher) GetDescriptorByWatch(wd int) *WatchDescriptor {
	w.Lock()
	defer w.Unlock()
	for _, d := range w.Descriptors {
		if d.WatchDescriptor == wd {
			return d
		}
	}
	return nil
}

// GetDescriptorByPath searches Watcher w for a watch descriptor.
// Searches by WatchDescriptor's path
func (w *Watcher) GetDescriptorByPath(watchPath string) *WatchDescriptor {
	if w.DescriptorExists(watchPath) == true {
		w.Lock()
		d := w.Descriptors[watchPath]
		w.Unlock()
		return d
	}
	return nil
}

func (w *Watcher) incrementEventCount() {
	atomic.AddUint32(&w.EventCount, 1)
}

// GetEventCount returns the atomic counter tracking the count of events for this Watcher. atomic/thread-safe.
func (w *Watcher) GetEventCount() uint32 {
	return atomic.LoadUint32(&w.EventCount)
}

// ReadSingleEvent reads and returns a single event from the watch descriptor.
func (w *Watcher) ReadSingleEvent() (*FsEvent, error) {
	var buffer [unix.SizeofInotifyEvent + unix.PathMax]byte

	bytesRead, err := C.read(C.int(w.InotifyDescriptor),
		unsafe.Pointer(&buffer),
		C.ulong(unix.SizeofInotifyEvent+unix.PathMax))

	if bytesRead < unix.SizeofInotifyEvent {
		return nil, ErrIncompleteRead
	} else if err != nil {
		return nil, fmt.Errorf("%s: %s", ErrReadError.Error(), err)
	}

	rawEvent := (*unix.InotifyEvent)(unsafe.Pointer(&buffer))

	descriptor := w.GetDescriptorByWatch(int(rawEvent.Wd))
	if descriptor == nil {
		return nil, ErrDescForEventNotFound
	}

	bytes := (*[unix.PathMax]byte)(unsafe.Pointer(&buffer[unix.SizeofInotifyEvent]))
	eventName := strings.TrimRight(string(bytes[0:rawEvent.Len]), "\000")
	eventPath := path.Clean(path.Join(descriptor.Path, eventName))

	event := &FsEvent{
		Name:       eventName,
		Path:       eventPath,
		Descriptor: descriptor,
		RawEvent:   rawEvent,
		ID:         w.GetEventCount(),
		Timestamp:  time.Now().UTC(),
	}
	w.incrementEventCount()
	return event, nil
}

// Watch calls ReadSingleEvent (which read-blocks) in a loop while there are running WatchDescriptors in Watcher w
// Writes events and errors to the channels w.Errors and w.Events
func (w *Watcher) Watch() {
	for w.GetRunningDescriptors() > 0 {
		event, err := w.ReadSingleEvent()
		if err != nil {
			w.Errors <- err
			continue
		}
		if event != nil {
			w.Events <- event
		}
	}
}

// RegisterEventHandler registers an EventHandler with Watcher.
// The EventHandler handle will be applied to any event read by WatchAndHandle matching its mask
func (w *Watcher) RegisterEventHandler(handle EventHandler) error {
	w.Lock()
	defer w.Unlock()

	for _, existingHandle := range w.eventHandlers {
		if existingHandle.GetMask() == handle.GetMask() {
			return ErrHandleExists
		}
	}
	w.eventHandlers = append(w.eventHandlers, handle)
	return nil
}

// UnregisterEventHandler Remove an EventHandle from a Watcher's EventHandle list
// The EventHandler handle will no longer applied to any event read by WatchAndHandle
func (w *Watcher) UnregisterEventHandler(removeMask uint32) error {
	w.Lock()
	defer w.Unlock()

	for index, handler := range w.eventHandlers {
		if handler.GetMask() == removeMask {
			if (index + 1) > len(w.eventHandlers) {
				w.eventHandlers = w.eventHandlers[:len(w.eventHandlers)-1]
			}
			w.eventHandlers = append(w.eventHandlers[:index], w.eventHandlers[index+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%s: event mask: %d", ErrNoSuchHandle, removeMask)
}

// getEventHandle returns the EventHandle matching event.RawEvent.Mask
func (w *Watcher) getEventHandle(event *FsEvent) EventHandler {
	w.Lock()
	defer w.Unlock()
	for _, handle := range w.eventHandlers {
		if handle.Check(event) == true {
			return handle
		}
	}
	return nil
}

/*
WatchAndHandle calls ReadSingleEvent to read an event, passing it to getEventHandle to retrieve
the correct handle for the event mask. Errors returned by the event's Handle are written to w.Errors
The event is *not* written to the w.Events channel.
If there is no handle registered to handle a specific event in the Watcher, WatchAndHandle immediately writes
ErrNoSuchHandle to the w.Errors channel and returns.
If there are no running watch descriptors, WatchAndHandle immediately writes ErrNoRunningDescriptors to w.Errors and returns.
If there are no registered EventHandles in the Watcher, WatchAndHandle immediately writes ErrNoEventHandles to w.Errors and returns.
*/
func (w *Watcher) WatchAndHandle() {
	for w.GetRunningDescriptors() > 0 && len(w.eventHandlers) > 0 {
		event, err := w.ReadSingleEvent()
		if err != nil {
			w.Errors <- err
			continue
		}
		if event != nil {
			if h := w.getEventHandle(event); h != nil {
				err := h.Handle(w, event)

				if err != nil {
					returnErr := errors.New(ErrHandleError.Error() + ": " + err.Error())
					w.Errors <- returnErr
				}
			} else {
				w.Errors <- fmt.Errorf("%s: event mask: %d", ErrNoSuchHandle, event.RawEvent.Mask)
			}
		}
	}

	if w.GetRunningDescriptors() == 0 {
		w.Errors <- ErrNoRunningDescriptors
	}
	if len(w.eventHandlers) == 0 {
		w.Errors <- ErrNoEventHandles
	}
}
