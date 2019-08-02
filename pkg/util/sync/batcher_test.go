package sync

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

var events = []string{"CREATE", "MODIFY", "DELETE"}

type job struct {
	file  string
	event string
}

func TestBatchWriter(t *testing.T) {
	ch := make(chan job)
	b := NewBatchWriter(1 * time.Second)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
			if i == 4 {
				fmt.Println("sleep")
				time.Sleep(2 * time.Second)
			}

			file := fmt.Sprintf("foo%d", i%5)
			event := events[i%3]

			eventList := []string{}
			val, ok := b.Load(file)
			if ok {
				eventList = val.([]string)
			}
			eventList = append(eventList, event)

			// Store and batch all existing changes after the timeout duration
			b.Store(file, eventList)
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
					file:  file,
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
