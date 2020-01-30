package dispatcher

import (
	"context"

	"github.com/go-loadtest/evbundler"
	"github.com/go-loadtest/evbundler/event"
)

type Dispatcher interface {
	Dispatch(context.Context, chan event.Event)
	Export() evbundler.Metrics
}
