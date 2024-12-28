package queue

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	ID           int
	Queue        *Queue
	EventHandler func(event *Event) error
}

type WorkerPool struct {
	workers    []*Worker
	numWorkers int
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int, queue *Queue, handler func(event *Event) error) *WorkerPool {
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

func (w *Worker) Start(done chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-done:
			// Exit when shutdown signal is received
			fmt.Printf("Worker %d: Shutdown signal received. Exiting. \n", w.ID)
			return
		default:
			event, err := w.Queue.Dequeue()
			if err != nil {
				// Queue is empty, wait for new events or shutdown signal
				time.Sleep(time.Second)
				continue
			}
			if err := w.EventHandler(event); err != nil {
				fmt.Printf("Worker %d failed to process event %s: %v\n", w.ID, event.ID, err)
			} else {
				fmt.Printf("Worker %d processed event %s\n", w.ID, event.ID)
			}
		}
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

func (wp *WorkerPool) Shutdown(done chan struct{}) {
	close(done)
	wp.Wait()
}
