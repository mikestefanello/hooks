package hooks

import (
	"testing"
)

func TestNewEvent(t *testing.T) {
	msg := &message{id: 100}
	h := NewHook[message](hookName)
	e := newEvent(h, msg)

	if e.Msg != msg {
		t.Fail()
	}

	if e.Hook != h {
		t.Fail()
	}
}
