package queue

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	ID           int
	Queue        *Queue
	EventHandler func(ctx context.Context, event *Event) error
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

			// Retry Logic using exponential backoff
			retryDelay := event.RetryDelay
			for event.RetryCount < event.MaxRetries {
				// Use event's timeout if specified
				ctx, cancel := context.WithTimeout(context.Background(), event.Timeout)
				defer cancel()

				if err := w.EventHandler(ctx, event); err != nil {
					event.RetryCount++
					// Apply exponential backoff
					time.Sleep(retryDelay)
					retryDelay *= 2
					if retryDelay > event.MaxDelay {
						retryDelay = event.MaxDelay
					}
					fmt.Printf("Worker %d failed to process event %s (Retry %d/%d) - Retrying in %v\n", w.ID, event.ID, event.RetryCount, event.MaxRetries, retryDelay)
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
		}
	}
}
