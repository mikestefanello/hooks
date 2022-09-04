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

		assertEqual(t, user, event.Msg, "user not equal to event message")
		assertEqual(t, testEventType, event.Type, "event type incorrect")
		assertEqual(t, testEventName, event.Type.GetName(), "event type name incorrect")
	}

	for i := 0; i < listenerCount; i++ {
		testEventType.Listen(listener)
	}

	testEventType.
		NewEvent(user).
		Dispatch()

	assertEqual(
		t,
		listenerCount,
		counter,
		"expected %d invocations of the listener - got %d",
		listenerCount,
		counter,
	)
}

func assertEqual(t *testing.T, expected, actual any, message string, args ...any) {
	if expected != actual {
		t.Errorf(message, args...)
	}
}
