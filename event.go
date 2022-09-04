package hooks

type Event[T any] struct {
	Msg  *T
	Hook *Hook[T]
}

// newEvent creates a new event
func newEvent[T any](hook *Hook[T], message *T) Event[T] {
	return Event[T]{
		Msg:  message,
		Hook: hook,
	}
}
