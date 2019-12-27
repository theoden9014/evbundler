package evbundler

import (
	"sync"

	"github.com/theoden9014/evbundler/event"
)

type RoundRobin struct {
	mu       sync.Mutex
	resource []interface{}
}

func (rr *RoundRobin) Schedule(ev event.Event) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	var e interface{}
	e, rr.resource = rr.resource[0], rr.resource[1:]
	rr.resource = append(rr.resource, e)
	// return e
}

type GoChannel struct {
	channel chan event.Event
}

func (gc GoChannel) Schedule(ev event.Event) {
	gc.channel <- ev
}
