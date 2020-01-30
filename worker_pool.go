package evbundler

import (
	"sync"
)

type WorkerPool interface {
	Get() Worker
	Put(Worker)
	Len() int
}

func NewWorkerPool(n int, f WorkerFunc) *workerPool {
	wp := &workerPool{}
	if f == nil {
		f = defaultWorkerFunc
	}

	for i := 0; i < n; i++ {
		w := &worker{f: f}
		wp.Put(w)
	}

	return wp
}

type workerPool struct {
	mu   sync.RWMutex
	pool []Worker
}

// Get returns a Worker from pool.
// return nil if empty in pool
func (wp *workerPool) Get() Worker {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if len(wp.pool) == 0 {
		return nil
	}

	var w Worker
	w, wp.pool = wp.pool[0], wp.pool[1:]

	return w
}

// Put store a Worker into pool.
func (wp *workerPool) Put(w Worker) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	wp.pool = append(wp.pool, w)
}

// Len returns count of workers in pool.
func (wp *workerPool) Len() int {
	return len(wp.pool)
}
