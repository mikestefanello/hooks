package hooks

import (
	"testing"
)

func TestHook_Dispatch(t *testing.T) {
	l, h, msg := newListener(t)

	h.Dispatch(msg)

	if listenerCount != l.counter {
		t.Fail()
	}
}

func TestHook_DispatchAsync(t *testing.T) {
	l, h, msg := newListener(t)

	h.DispatchAsync(msg)
	l.wg.Wait()

	if listenerCount != l.counter {
		t.Fail()
	}
}
