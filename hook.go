package hooks

import (
	"sync"
)

// Listener is a function that can listen and react to a hook event
type Listener[T any] func(event Event[T])

// Hook is a mechanism which supports the ability to dispatch data to arbitrary listener callbacks
type Hook[T any] struct {
	// name stores the name of the hook
	name string

	// listeners stores the functions which will be invoked during dispatch
	listeners []Listener[T]

	// mu stores the mutex to provide concurrency-safe operations
	mu sync.RWMutex
}

// NewHook creates a new Hook
func NewHook[T any](name string) *Hook[T] {
	logf("hook created: %s", name)

	return &Hook[T]{
		name:      name,
		listeners: make([]Listener[T], 0),
		mu:        sync.RWMutex{},
	}
}

// GetName returns the hook's name
func (h *Hook[T]) GetName() string {
	return h.name
}

// Listen registers a callback function to be invoked when the hook dispatches data
func (h *Hook[T]) Listen(callback Listener[T]) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.listeners = append(h.listeners, callback)

	logf("registered listener with hook: %s", h.GetName())
}

// GetListenerCount returns the number of listeners currently registered
func (h *Hook[T]) GetListenerCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.listeners)
}

// Dispatch invokes all listeners synchronously with the provided message
func (h *Hook[T]) Dispatch(message T) {
	h.dispatch(message, false)
}

// DispatchAsync invokes all listeners asynchronously with the provided message
func (h *Hook[T]) DispatchAsync(message T) {
	h.dispatch(message, true)
}

// dispatch invokes all listeners either synchronously or asynchronously with the provided message
func (h *Hook[T]) dispatch(message T, async bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	e := newEvent[T](h, message)

	// Check if the logger is available here to avoid the call since dispatching can happen very often and
	// this can help with performance
	if logger != nil {
		logf("dispatching hook %s to %d listeners (async: %v)", h.GetName(), len(h.listeners), async)
		defer logf("dispatch to hook %s complete", h.GetName())
	}

	for _, callback := range h.listeners {
		if async {
			go callback(e)
		} else {
			callback(e)
		}
	}
}
