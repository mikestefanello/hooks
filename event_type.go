package hooks

import (
	"sync"
)

type EventType[T any] struct {
	name      string
	listeners []func(event Event[T])
	mu        sync.RWMutex
}

// NewEventType creates a new event type
func NewEventType[T any](name string) *EventType[T] {
	return &EventType[T]{
		name:      name,
		listeners: make([]func(event Event[T]), 0),
		mu:        sync.RWMutex{},
	}
}

func (e *EventType[T]) GetName() string {
	return e.name
}

func (e *EventType[T]) Listen(callback func(event Event[T])) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.listeners = append(e.listeners, callback)
}

func (e *EventType[T]) GetListenerCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return len(e.listeners)
}

func (e *EventType[T]) NewEvent(message *T) *Event[T] {
	return NewEvent[T](e, message)
}
