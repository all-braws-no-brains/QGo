package tests

import (
	"main/pkg/queue"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool_GracefulShutdown(t *testing.T) {
	// Create the queue and worker pool
	q := queue.NewQueue() // q is now properly initialized using the constructor
	handler := func(event *queue.Event) error {
		// Simulate event processing (sleep for a short duration)
		time.Sleep(100 * time.Millisecond)
		return nil
	}
	pool := queue.NewWorkerPool(3, q, handler)

	// Start the worker pool with the shutdown channel
	done := make(chan struct{})
	pool.Start(done)

	// Enqueue events
	q.Enqueue(&queue.Event{ID: "1", Payload: []byte("event 1")})
	q.Enqueue(&queue.Event{ID: "2", Payload: []byte("event 2")})

	// Wait a bit for workers to process the events
	time.Sleep(500 * time.Millisecond)

	// Shut down the worker pool
	pool.Shutdown(done)

	// Test: Verify that the queue is empty after processing
	assert.True(t, q.IsEmpty(), "Queue should be empty after processing")
}
