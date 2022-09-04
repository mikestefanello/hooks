package hooks

import (
	"testing"
)

type testUser struct {
	id   int
	name string
}

const (
	testHookName  = "test.event.type"
	listenerCount = 3
)

func TestNewEvent(t *testing.T) {
	counter := 0
	user := &testUser{
		name: "Mike",
		id:   123,
	}

	testHook := NewHook[testUser](testHookName)

	listener := func(event Event[testUser]) {
		counter++

		if user != event.Msg {
			t.Fail()
		}

		if testHook != event.Hook {
			t.Fail()
		}

		if testHookName != event.Hook.GetName() {
			t.Fail()
		}
	}

	for i := 0; i < listenerCount; i++ {
		testHook.Listen(listener)
	}

	if listenerCount != testHook.GetListenerCount() {
		t.Fail()
	}

	testHook.Dispatch(user)

	if listenerCount != counter {
		t.Fail()
	}
}

//func thing() {
//	userPreUpdate := NewHook[testUser]("user.preupdate")
//	userUpdate := NewHook[testUser]("user.update")
//	userUpdate.Listen(func(event Event[testUser]) {
//
//	})
//
//	userPreUpdate.NewEvent(&testUser{}).Dispatch()
//	user.Save()
//	userUpdate.NewEvent(&testUser{}).Dispatch()
//}
