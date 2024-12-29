package tests

import (
	"testing"

	"github.com/all-braws-no-brains/QGo/pkg/queue"

	"github.com/stretchr/testify/assert"
)

func TestQueue_Enqueue_Dequeue(t *testing.T) {
	q := queue.NewQueue()

	// Test Enqueue
	event1 := &queue.Event{ID: "1", Payload: []byte("Hello")}
	q.Enqueue(event1)

	// Test Dequeue
	dequeuedEvent, err := q.Dequeue()
	assert.Nil(t, err, "Expected no error while dequeuing")
	assert.Equal(t, event1, dequeuedEvent, "Dequeued event should match the enqueued event")
}

func TestQueue_Dequeue_EmptyQueue(t *testing.T) {
	q := queue.NewQueue()

	// Test: Dequeue from an empty queue
	_, err := q.Dequeue()
	assert.NotNil(t, err, "Expected error when dequeuing from an empty queue")
}

func TestQueue_MultipleEvents(t *testing.T) {
	q := queue.NewQueue()

	// Test: Enqueue multiple events
	event1 := &queue.Event{ID: "1", Payload: []byte("First")}
	event2 := &queue.Event{ID: "2", Payload: []byte("Second")}
	q.Enqueue(event1)
	q.Enqueue(event2)

	// Test: Dequeue events in the same order
	dequeuedEvent1, err1 := q.Dequeue()
	dequeuedEvent2, err2 := q.Dequeue()

	assert.Nil(t, err1, "Expected no error when dequeuing")
	assert.Nil(t, err2, "Expected no error when dequeuing")
	assert.Equal(t, event1, dequeuedEvent1, "Dequeued first event should match event1")
	assert.Equal(t, event2, dequeuedEvent2, "Dequeued second event should match event2")
}
