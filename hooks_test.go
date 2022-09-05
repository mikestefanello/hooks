package hooks

import (
	"sync"
	"testing"
)

const (
	hookName      = "test.hook"
	listenerCount = 3
)

type message struct {
	id int
}

type listener struct {
	counter int
	msg     *message
	hook    *Hook[message]
	wg      sync.WaitGroup
	t       *testing.T
}

// newListener creates and initializes a new listener with a hook and message
func newListener(t *testing.T) (*listener, *Hook[message], *message) {
	msg := &message{id: 123}
	h := NewHook[message](hookName)
	l := &listener{
		t:    t,
		msg:  msg,
		hook: h,
	}

	for i := 0; i < listenerCount; i++ {
		h.Listen(l.Callback)
	}

	l.wg.Add(listenerCount)

	if listenerCount != h.GetListenerCount() {
		t.Fail()
	}

	return l, h, msg
}

func (l *listener) Callback(event Event[message]) {
	l.counter++

	if l.msg != event.Msg {
		l.t.Fail()
	}

	if l.hook != event.Hook {
		l.t.Fail()
	}

	if hookName != event.Hook.GetName() {
		l.t.Fail()
	}

	l.wg.Done()
}
