package dispatcher

import (
	"context"
	"runtime/trace"

	"github.com/theoden9014/evbundler"
)

type RoundRobin struct {
	pool     evbundler.WorkerPool
	resultCh chan *evbundler.Result
	metrics  evbundler.Metrics
}

func (d *RoundRobin) Dispatch(ctx context.Context, evCh chan evbundler.Event) error {
	go d.receiveResult(ctx)
	for {
		select {
		case ev := <-evCh:
			go func() { _ = d.dispatch(ctx, ev) }()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (d *RoundRobin) dispatch(ctx context.Context, ev evbundler.Event) error {
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
	region.End()

	return nil
}

func (d *RoundRobin) receiveResult(ctx context.Context) {
	for {
		select {
		case r := <-d.resultCh:
			d.metrics.Add(r)
		case <-ctx.Done():
			return
		}
	}
}
