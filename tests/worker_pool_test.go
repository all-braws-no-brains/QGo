package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/all-braws-no-brains/QGo/pkg/queue"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool_GracefulShutdown(t *testing.T) {
	// Create the queue and worker pool
	q := queue.NewQueue() // q is now properly initialized using the constructor
	handler := func(ctx context.Context, event *queue.Event) error {
		// Simulate event processing (sleep for a short duration)
		time.Sleep(100 * time.Millisecond)
		return nil
	}
	pool := queue.NewWorkerPool(3, q, handler)

	// Start the worker pool with the shutdown channel
	done := make(chan struct{})
	pool.Start(done)

	// Enqueue events
	q.Enqueue(&queue.Event{
		ID:         "1",
		Payload:    []byte("event 1"),
		RetryCount: 0, // No retries
		MaxRetries: 3,
		RetryDelay: 100 * time.Millisecond,
	})
	q.Enqueue(&queue.Event{
		ID:         "2",
		Payload:    []byte("event 2"),
		RetryCount: 0, // No retries
		MaxRetries: 3,
		RetryDelay: 100 * time.Millisecond,
	})

	// Wait a bit for workers to process the events
	time.Sleep(500 * time.Millisecond)

	// Shut down the worker pool
	pool.Shutdown(done)

	// Test: Verify that the queue is empty after processing
	assert.True(t, q.IsEmpty(), "Queue should be empty after processing")
}

func TestWorkerPool_Concurrency(t *testing.T) {
	// Create a new queue
	q := queue.NewQueue()

	// Define an event handler that just prints the event
	handler := func(ctx context.Context, event *queue.Event) error {
		fmt.Printf("Processing event: %s\n", event.ID)
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	// Create a worker pool with 3 workers
	wp := queue.NewWorkerPool(3, q, handler)

	// Create a done channel for graceful shutdown
	done := make(chan struct{})

	// Start the worker pool
	wp.Start(done)

	// Enqueue 10 events into the queue
	for i := 0; i < 10; i++ {
		event := &queue.Event{
			ID:         fmt.Sprintf("event-%d", i),
			Payload:    []byte(fmt.Sprintf("Payload of event %d", i)),
			RetryCount: 0,
			MaxRetries: 3,
			RetryDelay: 100 * time.Millisecond,
		}
		q.Enqueue(event)
	}

	// Allow workers some time to process events
	time.Sleep(2 * time.Second)

	// Shutdown the worker pool gracefully
	wp.Shutdown(done)

	// Wait for workers to finish
	wp.Wait()

	// Check that all events have been processed
	// (This could be enhanced by tracking which events were processed)
	t.Log("All events have been processed")
}

func TestWorkerPool_EventRetry(t *testing.T) {
	q := queue.NewQueue()

	// Define a handler that simulates failure for the first 2 attempts
	handler := func(ctx context.Context, event *queue.Event) error {
		if event.RetryCount < 2 {
			return fmt.Errorf("temporary failure")
		}
		return nil
	}

	// Create a worker pool
	pool := queue.NewWorkerPool(3, q, handler)

	// Start the worker pool with shutdown channel
	done := make(chan struct{})
	pool.Start(done)

	// Enqueue an event that will fail twice
	event := &queue.Event{ID: "1", Payload: []byte("Test Event"), MaxRetries: 3, RetryDelay: 500 * time.Millisecond}
	q.Enqueue(event)

	// Allow workers some time to process
	time.Sleep(3 * time.Second)

	// Shutdown the worker pool gracefully
	pool.Shutdown(done)

	// Wait for workers to finish
	pool.Wait()

	// Verify that the event has been processed after retries
	assert.True(t, q.IsEmpty(), "Queue should be empty after processing")
}
