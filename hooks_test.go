package hooks

import (
	"testing"
)

type testUser struct {
	id   int
	name string
}

const (
	testEventName = "test.event.type"
	listenerCount = 3
)

func TestNewEvent(t *testing.T) {
	counter := 0
	user := &testUser{
		name: "Mike",
		id:   123,
	}

	testEventType := NewEventType[testUser](testEventName)

	listener := func(event Event[testUser]) {
		counter++

		if user != event.Msg {
			t.Fail()
		}

		if testEventType != event.Type {
			t.Fail()
		}

		if testEventName != event.Type.GetName() {
			t.Fail()
		}
	}

	for i := 0; i < listenerCount; i++ {
		testEventType.Listen(listener)
	}

	if listenerCount != testEventType.GetListenerCount() {
		t.Fail()
	}

	testEventType.
		NewEvent(user).
		Dispatch()

	if listenerCount != counter {
		t.Fail()
	}
}
