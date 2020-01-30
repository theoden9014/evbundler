package evbundler

import (
	"context"
	"fmt"
	"runtime/trace"
	"time"

	"github.com/go-loadtest/evbundler/event"
	"golang.org/x/time/rate"
)

type Worker interface {
	Process(context.Context, event.Event) *Result
	Close() error
	// StateTransaction(<-chan WorkerState)
}

type worker struct {
	state WorkerState
	f     WorkerFunc
	limit *rate.Limiter
	r     <-chan event.Event
}

func (w *worker) Wait(ctx context.Context) event.Event {
	w.setState(StateWaiting)
	defer w.setState(StateActive)
	_ = w.limit.Wait(ctx)
	ev := <-w.r

	return ev
}

func (w *worker) Process(ctx context.Context, ev event.Event) *Result {
	w.setState(StateProcess)
	defer w.setState(StateActive)
	defer trace.StartRegion(ctx, fmt.Sprintf("process event: %q", ev.Name())).End()

	start := time.Now()
	weight, err := w.f(ctx, ev)
	elapsed := time.Now().Sub(start)

	return &Result{
		Weight:    weight,
		EventName: ev.Name(),
		Error:     err,
		Latency:   elapsed,
		Timestamp: start,
	}
}

type WorkerFunc func(context.Context, event.Event) (int, error)

type WorkerState int

const (
	StateDead WorkerState = iota
	StateActive
	StateWaiting
	StateProcess
)

func (ws WorkerState) String() string {
	return []string{"DEAD", "ACTIVE", "WAIT_EVENT", "PROCESS_EVENT"}[ws]
}

func (w *worker) setState(state WorkerState) {
	w.state = state
}
