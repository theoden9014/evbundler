package evbundler

import (
	"context"
	"fmt"
	"runtime/trace"
	"sync"
	"time"

	"github.com/go-loadtest/evbundler/event"
	"golang.org/x/time/rate"
)

type WorkerPoolMetrics struct {
	Workers uint64 `json:"workers"`
	Waiting uint64 `json:"waiting_workers"`
	Process uint64 `json:"processing_workers"`
}

type workerPool struct {
	mu   sync.RWMutex
	pool []*worker
}

func (wp *workerPool) Iterator() <-chan *worker {
	ch := make(chan *worker)
	go func() {
		defer close(ch)
		for _, w := range wp.pool {
			ch <- w
		}
	}()

	return ch
}

func (wp *workerPool) Push(w *worker) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	wp.pool = append(wp.pool, w)
}

func (wp *workerPool) Len() int {
	return len(wp.pool)
}

func (wp *workerPool) Count() int {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	var c int
	for _, w := range wp.pool {
		if !w.IsDead() {
			c++
		}
	}

	return c
}

func (wp *workerPool) ProccessCount() int {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	var c int
	for _, w := range wp.pool {
		if w.IsProccess() {
			c++
		}
	}

	return c
}

func (wp *workerPool) WaitingCount() int {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	var c int
	for _, w := range wp.pool {
		if w.IsWaiting() {
			c++
		}
	}

	return c
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

type worker struct {
	state WorkerState
	id    string
	f     WorkerFunc
	limit *rate.Limiter
	r     <-chan event.Event
}

func (w *worker) setState(state WorkerState) {
	w.state = state
}

func (w *worker) WaitEvent(ctx context.Context) event.Event {
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
		EventName: w.id,
		Error:     err,
		Latency:   elapsed,
		Timestamp: start,
	}
}

func (w *worker) Close() error {
	w.setState(StateDead)
	return nil
}
func (w *worker) IsWaiting() bool {
	return w.state == StateWaiting
}

func (w *worker) IsProccess() bool {
	return w.state == StateProcess
}

func (w *worker) IsDead() bool {
	return w.state == StateDead
}

type WorkerFunc func(context.Context, event.Event) (int, error)
