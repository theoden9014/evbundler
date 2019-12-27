package evbundler

import (
	"context"

	"github.com/theoden9014/evbundler/event"
)

type EventBundler interface {
	Register(chan event.Event)
	Bundle() <-chan event.Event
	Start(context.Context)
}

// type eventBundler struct {
// 	mu      sync.RWMutex
// 	inputs  []<-chan Event
// 	outputs chan Event
// }
