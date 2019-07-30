package watch

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func NewBatcher(syncMap *sync.Map, duration time.Duration) *Batcher {
	return &Batcher{
		duration: duration,
		flushCh:  make(chan struct{}),
		syncMap:  syncMap,
	}
}

type Batcher struct {
	duration time.Duration
	timer    *time.Timer
	flushCh  chan struct{}
	syncMap  *sync.Map
}

func (b *Batcher) CancelUnfiredTimer() {
	// If the timer already exists; stop it
	if b.timer != nil {
		log.Tracef("Batcher: Cancelled timer")
		b.timer.Stop()
	}
}

func (b *Batcher) DispatchAfterTimeout() {
	b.timer = time.AfterFunc(b.duration, func() {
		log.Tracef("Batcher: Dispatching a batch job")
		b.flushCh <- struct{}{}
	})
}

func (b *Batcher) Close() {
	log.Tracef("Batcher: Closing the batch channel")
	close(b.flushCh)
}

func (b *Batcher) ProcessBatch(fn func(key, val interface{}) bool) bool {
	_, ok := <-b.flushCh
	if !ok {
		// channel is closed
		return false
	}
	b.syncMap.Range(fn)
	*b.syncMap = sync.Map{}
	return true
}
