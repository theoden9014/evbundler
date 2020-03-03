package evbundler

import (
	"context"
	"sync"
)

type Event interface {
	Name() string
	Fire(ctx context.Context) error
}

type Producer interface {
	Produce(ctx context.Context) chan Event
}

type Dispatcher interface {
	Dispatch(context.Context, chan Event) error
}

type WorkerPool interface {
	Get() Worker
	Put(Worker)
	Len() int
}

type Metrics interface {
	Add(*Result)
	MarshalJSON() ([]byte, error)
}

// EventBundler put several receive channels to gather into one.
type EventBundler struct {
	mu     sync.RWMutex
	inputs []<-chan Event
	output chan Event
}

// In register a receive channel.
func (ep *EventBundler) In(ev chan Event) {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	ep.inputs = append(ep.inputs, ev)
}

// Out returns the bundled channels
func (ep *EventBundler) Out() <-chan Event {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	if ep.output == nil {
		ep.output = make(chan Event, len(ep.inputs))
	}
	return ep.output
}

// Start bypass receive channels to bundled channels.
// run another goroutine each receive channels.
func (ep *EventBundler) Start(ctx context.Context) {
	for _, evc := range ep.inputs {
		evc := evc
		go func(evc <-chan Event) {
			for {
				select {
				case ev := <-evc:
					ep.output <- ev
				case <-ctx.Done():
					return
				}
			}
		}(evc)
	}
}
