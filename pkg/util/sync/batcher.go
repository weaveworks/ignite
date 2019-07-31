package sync

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// NewBatchWriter creates a new BatchWriter
func NewBatchWriter(duration time.Duration) *BatchWriter {
	return &BatchWriter{
		duration: duration,
		flushCh:  make(chan struct{}),
		syncMap:  &sync.Map{},
	}
}

// BatchWriter is a struct that wraps a concurrent sync.Map
// and dispatches all writes to it at once, a specific
// duration (e.g. 1s) after the last write was performed. This
// allows for 100s of concurrent writes in milliseconds to the
// same map in one sending goroutine; and one receiving goroutine
// which can process the result after all the writes are done.
type BatchWriter struct {
	duration time.Duration
	timer    *time.Timer
	flushCh  chan struct{}
	syncMap  *sync.Map
}

// Load reads the key from the map
func (b *BatchWriter) Load(key interface{}) (value interface{}, ok bool) {
	return b.syncMap.Load(key)
}

// Store writes the value for the specified key to the map
// If no other .Store call is made during the specified duration,
// flushCh is invoked and ProcessBatch unblocks in the other goroutine
func (b *BatchWriter) Store(key, value interface{}) {
	// prevent the timer from firing as we're manipulating it now
	b.cancelUnfiredTimer()
	// store the key and the value as requested
	log.Tracef("BatchWriter: Storing key %v and value %q, reset the timer.", key, value)
	b.syncMap.Store(key, value)
	// set the timer to fire after the duration, unless there's a new .Store call
	b.dispatchAfterTimeout()
}

// Close closes the underlying channel
func (b *BatchWriter) Close() {
	log.Trace("BatchWriter: Closing the batch channel")
	close(b.flushCh)
}

// ProcessBatch is effectively a Range over the sync.Map, once a batch write is
// released. This should be used in the receiving goroutine. The internal map is
// reset after this call, so be sure to capture all the contents if needed. This
// function returns false if Close() has been called.
func (b *BatchWriter) ProcessBatch(fn func(key, val interface{}) bool) bool {
	if _, ok := <-b.flushCh; !ok {
		// channel is closed
		return false
	}
	log.Trace("BatchWriter: Received a flush for the batch. Dispatching it now.")
	b.syncMap.Range(fn)
	*b.syncMap = sync.Map{}
	return true
}

func (b *BatchWriter) cancelUnfiredTimer() {
	// If the timer already exists; stop it
	if b.timer != nil {
		log.Tracef("BatchWriter: Cancelled timer")
		b.timer.Stop()
		b.timer = nil
	}
}

func (b *BatchWriter) dispatchAfterTimeout() {
	b.timer = time.AfterFunc(b.duration, func() {
		log.Tracef("BatchWriter: Dispatching a batch job")
		b.flushCh <- struct{}{}
	})
}
