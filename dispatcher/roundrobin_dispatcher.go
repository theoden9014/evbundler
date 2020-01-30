package dispatcher

import (
	"context"
	"runtime/trace"

	"github.com/go-loadtest/evbundler"
	"github.com/go-loadtest/evbundler/event"
)

type RoundRobinDispatcher struct {
	pool     evbundler.WorkerPool
	resultCh chan *evbundler.Result
	metrics  *evbundler.Metrics
}

func (d *RoundRobinDispatcher) Dispatch(ctx context.Context, evCh chan event.Event) error {
	for {
		select {
		case ev := <-evCh:
			go d.dispatch(ctx, ev)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (d *RoundRobinDispatcher) dispatch(ctx context.Context, ev event.Event) error {
	w := d.pool.Get()
	defer d.pool.Put(w)

	ctx, task := trace.NewTask(ctx, "worker")
	defer task.End()
	if ev == nil { // close Event channel
		return nil
	}

	r := w.Process(ctx, ev)
	d.resultCh <- r

	region := trace.StartRegion(ctx, "send result")
	d.resultCh <- r
	region.End()

	return nil
}
