package main

import (
	"context"
	"fmt"

	"time"

	"github.com/all-braws-no-brains/QGo/pkg/queue"
	"github.com/all-braws-no-brains/QGo/pkg/utils"
)

// EmailHandler is an example handler for processing email-related events.
func EmailHandler(ctx context.Context, event *queue.Event, config *utils.EmailConfig) error {
	fmt.Printf("Processing email event: %s\n", event.ID)
	return queue.EmailEventHandler(ctx, event, config)
}

func main() {
	// Create a new queue
	q := queue.NewQueue()

	// Create an email configuration
	emailConfig := &utils.EmailConfig{
		SMTPHost: "smtp.example.com",
		Port:     "587",
		From:     "your-email@example.com",
		To:       "recipient@example.com",
		Subject:  "Order Confirmation",
	}

	// Create and enqueue some events
	event1 := &queue.Event{
		ID:         "email-event-1",
		Payload:    "Thank you for your order!",
		MaxRetries: 3,
		RetryDelay: 1 * time.Second,
		Timeout:    5 * time.Second,
		MaxDelay:   10 * time.Second,
	}

	q.Enqueue(event1)

	// Create a worker pool with the email handler
	workerPool := queue.NewWorkerPool(3, q, func(ctx context.Context, event *queue.Event) error {
		return EmailHandler(ctx, event, emailConfig)
	})

	// Channel to signal workers to stop
	done := make(chan struct{})

	// Start the worker pool
	workerPool.Start(done)

	// Allow workers to process events for 10 seconds
	time.Sleep(10 * time.Second)

	// Shut down the worker pool
	workerPool.Shutdown(done)
	fmt.Println("All workers have stopped.")
}
