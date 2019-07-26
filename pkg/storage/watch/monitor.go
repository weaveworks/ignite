package watch

import "sync"

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
