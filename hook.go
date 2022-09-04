package hooks

import (
	"sync"
)

type Hook[T any] struct {
	name      string
	listeners []func(event Event[T])
	mu        sync.RWMutex
}

// NewHook creates a new hook
func NewHook[T any](name string) *Hook[T] {
	return &Hook[T]{
		name:      name,
		listeners: make([]func(event Event[T]), 0),
		mu:        sync.RWMutex{},
	}
}

func (h *Hook[T]) GetName() string {
	return h.name
}

func (h *Hook[T]) Listen(callback func(event Event[T])) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.listeners = append(h.listeners, callback)
}

func (h *Hook[T]) GetListenerCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.listeners)
}

func (h *Hook[T]) Dispatch(message *T) {
	h.dispatch(message, false)
}

func (h *Hook[T]) DispatchAsync(message *T) {
	h.dispatch(message, true)
}

func (h *Hook[T]) dispatch(message *T, async bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	e := newEvent[T](h, message)

	for _, callback := range h.listeners {
		if async {
			go callback(e)
		} else {
			callback(e)
		}
	}
}
