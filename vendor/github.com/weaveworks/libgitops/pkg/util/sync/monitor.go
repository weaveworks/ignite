package sync

import "sync"

// Monitor is a convenience wrapper around
// starting a goroutine with a wait group,
// which can be used to wait for the
// goroutine to stop.
type Monitor struct {
	wg *sync.WaitGroup
}

func RunMonitor(f func()) (m *Monitor) {
	m = &Monitor{
		wg: new(sync.WaitGroup),
	}

	m.wg.Add(1)
	go func() {
		f()
		m.wg.Done()
	}()

	return
}

func (m *Monitor) Wait() {
	if m != nil {
		m.wg.Wait()
	}
}
