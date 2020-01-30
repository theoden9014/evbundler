package evbundler

import (
	"context"
	"fmt"
	"runtime/trace"
	"time"

	"github.com/go-loadtest/evbundler/event"
)

type Worker interface {
	Process(context.Context, event.Event) *Result
	Close() error
	// StateTransaction(<-chan WorkerState)
}

type WorkerFunc func(context.Context, event.Event) error

func defaultWorkerFunc(ctx context.Context, ev event.Event) error {
	return ev.Fire(ctx)
}

type worker struct {
	state WorkerState
	f     WorkerFunc
}

func (w *worker) Process(ctx context.Context, ev event.Event) *Result {
	w.setState(StateProcess)
	defer w.setState(StateActive)
	defer trace.StartRegion(ctx, fmt.Sprintf("process event: %q", ev.Name())).End()

	start := time.Now()
	err := w.f(ctx, ev)
	elapsed := time.Now().Sub(start)

	return &Result{
		Weight:    1,
		EventName: ev.Name(),
		Error:     err,
		Latency:   elapsed,
		Timestamp: start,
	}
}

func (w worker) Close() error {
	return nil
}

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
