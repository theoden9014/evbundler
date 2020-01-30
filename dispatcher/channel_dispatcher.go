package dispatcher

import (
	"context"
	"runtime/trace"

	"github.com/go-loadtest/evbundler"
	"github.com/go-loadtest/evbundler/event"
)

// ChannelDispatcher dispatcher by go channel.
type ChannelDispatcher struct {
	pool     evbundler.WorkerPool
	resultCh chan *evbundler.Result
	metrics  *evbundler.Metrics
}

// NewChannelDispatcher initialize ChannelDispatcher.
func NewChannelDispatcher(pool evbundler.WorkerPool) *ChannelDispatcher {
	d := &ChannelDispatcher{
		pool:     pool,
		resultCh: make(chan *evbundler.Result, pool.Len()),
		metrics:  &evbundler.Metrics{},
	}
	return d
}

// Dispatch dispatches events from a event channel.
func (d *ChannelDispatcher) Dispatch(ctx context.Context, evCh chan event.Event) {
	go d.startWorkers(ctx, evCh)
	d.receiveResult(ctx)
}

func (d *ChannelDispatcher) receiveResult(ctx context.Context) {
	for {
		select {
		case r := <-d.resultCh:
			d.metrics.Add(r)
		case <-ctx.Done():
			return
		}
	}
}

func (d *ChannelDispatcher) startWorkers(ctx context.Context, evCh chan event.Event) {
	for w := d.pool.Get(); w != nil; w = d.pool.Get() {
		w := w
		go func() {
			defer w.Close()
			for {
				select {
				case ev := <-evCh:
					func() {
						ctx, task := trace.NewTask(ctx, "worker")
						defer task.End()
						if ev == nil { // close Event channel
							return
						}
						r := w.Process(ctx, ev)

						region := trace.StartRegion(ctx, "send result")
						d.resultCh <- r
						region.End()
					}()

				case <-ctx.Done():
					return
				}
			}
		}()
	}
}
