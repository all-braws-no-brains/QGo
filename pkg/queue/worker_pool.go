package queue

import (
	"context"
	"fmt"
	"sync"
)

type WorkerPool struct {
	workers    []*Worker
	numWorkers int
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int, queue *Queue, handler func(ctx context.Context, event *Event) error) *WorkerPool {
	pool := WorkerPool{
		numWorkers: numWorkers,
		workers:    make([]*Worker, numWorkers),
	}

	for i := 0; i < numWorkers; i++ {
		pool.workers[i] = &Worker{
			ID:           i + 1,
			Queue:        queue,
			EventHandler: handler,
		}
	}
	return &pool
}

func (wp *WorkerPool) Start(done chan struct{}) {
	for _, worker := range wp.workers {
		wp.wg.Add(1)
		go worker.Start(done, &wp.wg)
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

func (wp *WorkerPool) Shutdown(done chan struct{}) {
	fmt.Println("Shutting down worker pool...")
	close(done) // Signal workers to stop
	wp.Wait()   // Wait for all workers to finish
	fmt.Println("Worker pool has been shut down.")
}
