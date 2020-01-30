package evbundler

import (
	"sync"
)

type WorkerPool interface {
	Get() Worker
	Put(Worker)
	Len() int
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

// Push restore a Worker into pool.
func (wp *workerPool) Push(w Worker) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	wp.pool = append(wp.pool, w)
}

// Len returns count of workers in pool.
func (wp *workerPool) Len() int {
	return len(wp.pool)
}
