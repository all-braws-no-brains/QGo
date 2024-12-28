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

			// Retry Logic
			for event.RetryCount < event.MaxRetries {
				if err := w.EventHandler(event); err != nil {
					event.RetryCount++
					// Apply delay or back off
					time.Sleep(event.RetryDelay) // You can replace this with exponential backoff if needed
					fmt.Printf("Worker %d failed to process event %s (Retry %d/%d)\n", w.ID, event.ID, event.RetryCount, event.MaxRetries)
					continue
				}
				// Successfully processed the event
				fmt.Printf("Worker %d processed event %s\n", w.ID, event.ID)
				break
			}

			if event.RetryCount >= event.MaxRetries {
				// If max retries reached, log a failure
				fmt.Printf("Worker %d failed to process event %s after %d retries\n", w.ID, event.ID, event.RetryCount)
			}

			switch payload := event.Payload.(type) {
			case []byte:
				// If the payload is of type []byte, process it accordingly
				fmt.Printf("Worker %d processing payload: %s\n", w.ID, string(payload))
			default:
				// Handle other types or log them as unknown
				fmt.Printf("Worker %d processing unknown payload type: %T\n", w.ID, payload)
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
