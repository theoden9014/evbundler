package dispatcher

import (
	"context"
	"errors"
	"runtime/trace"
	"sync"

	"github.com/go-loadtest/evbundler"
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
func (d *ChannelDispatcher) Dispatch(ctx context.Context, evCh chan evbundler.Event) error {
	go d.receiveResult(ctx)
	return d.dispatch(ctx, evCh)
}

func (d *ChannelDispatcher) dispatch(ctx context.Context, evCh chan evbundler.Event) error {
	if d.pool.Len() == 0 {
		return errors.New("count of workers > 0 in worker pool")
	}

	var wg sync.WaitGroup
	for w := d.pool.Get(); w != nil; w = d.pool.Get() {
		w := w
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.serveWorker(ctx, w, evCh)
		}()
	}

	<-ctx.Done()
	// wait for worker graceful shutdown
	wg.Wait()
	return ctx.Err()
}

func (d *ChannelDispatcher) serveWorker(ctx context.Context, w evbundler.Worker, evCh chan evbundler.Event) {
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
