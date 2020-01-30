package evbundler

import (
	"context"
	"sync"

	"github.com/go-loadtest/evbundler/event"
)

// EventBundler put several receive channels to gather into one.
type EventBundler struct {
	mu     sync.RWMutex
	inputs []<-chan event.Event
	output chan event.Event
}

// In register a receive channel.
func (ep *EventBundler) In(ev chan event.Event) {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	ep.inputs = append(ep.inputs, ev)
}

// Out returns the bundled channels
func (ep *EventBundler) Out() <-chan event.Event {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	if ep.output == nil {
		ep.output = make(chan event.Event, len(ep.inputs))
	}
	return ep.output
}

// Start bypass receive channels to bundled channels.
// run another goroutine each receive channels.
func (ep *EventBundler) Start(ctx context.Context) {
	for _, evc := range ep.inputs {
		evc := evc
		go func(evc <-chan event.Event) {
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
