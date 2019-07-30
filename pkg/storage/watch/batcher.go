package watch

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func NewBatcher(duration time.Duration) *Batcher {
	return &Batcher{
		duration: duration,
		flushCh:  make(chan struct{}),
		syncMap:  &sync.Map{},
	}
}

type Batcher struct {
	duration time.Duration
	timer    *time.Timer
	flushCh  chan struct{}
	syncMap  *sync.Map
}

func (b *Batcher) Load(key interface{}) (value interface{}, ok bool) {
	return b.syncMap.Load(key)
}

func (b *Batcher) Store(key, value interface{}) {
	// prevent the timer from firing as we're manipulating it now
	b.cancelUnfiredTimer()
	// store the key and the value as requested
	log.Tracef("Batcher: Storing key %v and value %q, reset the timer.", key, value)
	b.syncMap.Store(key, value)
	// set the timer to fire after the duration, unless there's a new .Store call
	b.dispatchAfterTimeout()
}

func (b *Batcher) Close() {
	log.Trace("Batcher: Closing the batch channel")
	close(b.flushCh)
}

func (b *Batcher) ProcessBatch(fn func(key, val interface{}) bool) bool {
	_, ok := <-b.flushCh
	if !ok {
		// channel is closed
		return false
	}
	log.Trace("Batcher: Received a flush for the batch. Dispatching it now.")
	b.syncMap.Range(fn)
	*b.syncMap = sync.Map{}
	return true
}

func (b *Batcher) cancelUnfiredTimer() {
	// If the timer already exists; stop it
	if b.timer != nil {
		log.Tracef("Batcher: Cancelled timer")
		b.timer.Stop()
		b.timer = nil
	}
}

func (b *Batcher) dispatchAfterTimeout() {
	b.timer = time.AfterFunc(b.duration, func() {
		log.Tracef("Batcher: Dispatching a batch job")
		b.flushCh <- struct{}{}
	})
}
