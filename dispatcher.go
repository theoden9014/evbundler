package evbundler

import (
	"sort"

	"github.com/go-loadtest/evbundler/event"
)

type Worker interface{}

type WorkerPool interface {
	sort.Interface
}

type Dispatcher interface {
	Dispatch(event.Event)
}

type Scheduler interface {
	Schedule(event.Event) interface{}
}

type dispatcher struct {
	scheduler Scheduler
}

func (d dispatcher) Dispatch(ev event.Event) {
	d.scheduler.Schedule(ev)
}

// type channelDispatcher struct {
// 	mu     sync.Mutex
// 	inputs []<-chan Event
// 	output chan Event
// }

// func (d *channelDispatcher) Register(evch chan Event) {
// 	d.mu.Lock()
// 	defer d.mu.Unlock()

// 	d.inputs = append(d.inputs, evch)
// }

// func (d *channelDispatcher) Bundle()
