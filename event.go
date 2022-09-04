package hooks

type Event[T any] struct {
	Msg  *T
	Type *EventType[T]
}

// NewEvent creates a new event
func NewEvent[T any](eventType *EventType[T], message *T) *Event[T] {
	return &Event[T]{
		Msg:  message,
		Type: eventType,
	}
}

func (e *Event[T]) Dispatch() {
	e.dispatch(false)
}

func (e *Event[T]) DispatchAsync() {
	e.dispatch(true)
}

func (e *Event[T]) dispatch(async bool) {
	e.Type.mu.RLock()
	defer e.Type.mu.RUnlock()

	for _, callback := range e.Type.listeners {
		if async {
			go callback(*e)
		} else {
			callback(*e)
		}
	}
}
