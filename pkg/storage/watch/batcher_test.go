package watch

import (
	"testing"
	"sync"
	"fmt"
	"time"
	"strings"
)

var events = []string{"CREATE", "MODIFY", "DELETE"}

type job struct {
	file string
	event string
}

func TestBatcher(t *testing.T) {
	ch := make(chan job)
	syncMap := &sync.Map{} // the map is map[string][]string
	b := NewBatcher(syncMap, 1 * time.Second)
	go func() {
		var eventList []string
		for i := 0; i < 10; i++ {
			fmt.Println(i)
			if i == 4 {
				fmt.Println("sleep")
				time.Sleep(2 * time.Second)
			}
			// Notify the Batcher we got a new job
			b.CancelUnfiredTimer()

			file := fmt.Sprintf("foo%d", i % 5)
			event := events[i % 3]

			val, ok := syncMap.Load(file)
			if !ok {
				eventList = []string{event}
			} else {
				eventList = val.([]string)
				eventList = append(eventList, event)
			}

			syncMap.Store(file, eventList)

			// Batch all existing changes after the timeout duration
			b.DispatchAfterTimeout()
		}
		time.Sleep(2 * time.Second)
		fmt.Println("stopping")
		b.Close()
		close(ch)
	}()
	go func() {
		for {
			// When the batch items are available, process the items from the map
			ok := b.ProcessBatch(func(key, val interface{}) bool {
				file := key.(string)
				events := val.([]string)
				ch <- job{
					file: file,
					event: strings.Join(events, ","),
				}
				return true
			})
			if !ok {
				return
			}
			fmt.Println("")
			fmt.Println("Map flushed")
			fmt.Println("")
		}
	}()
	
	for j := range ch {
		fmt.Println(j.file, j.event)
	}
	
	//t.Error("err")
}