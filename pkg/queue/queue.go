package queue

import (
	"errors"
	"sync"
	"time"
)

type Event struct {
	ID         string
	Payload    interface{}
	RetryCount int
	MaxRetries int
	RetryDelay time.Duration
	// Timeout    time.Duration
}

type Queue struct {
	events []*Event
	mu     sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		events: make([]*Event, 0),
	}
}

func (q *Queue) Enqueue(event *Event) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.events = append(q.events, event)
	return nil
}

func (q *Queue) Dequeue() (*Event, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.events) == 0 {
		return nil, errors.New("queue is empty")
	}

	event := q.events[0]
	q.events = q.events[1:]
	return event, nil
}

func (q *Queue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.events) == 0
}

func (q *Queue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.events)
}
