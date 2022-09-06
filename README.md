# Hooks

## Overview

_Hooks_ provides a simple, type-safe hook system to enable easier modularization of your Go code. A _hook_ allows various parts of your codebase to tap in to events and operations happening elsewhere which prevents direct coupling between the producer and the consumers/listeners. For example, a _user_ package/module in your code may dispatch a _hook_ when a user is created, allowing your _notification_ package to send the user an email, and a _history_ package to record the activity without the _user_ module having to call these components directly. A hook can also be used to allow other modules to alter and extend data before it is processed.

## Installation

`go get github.com/mikestefanello/hooks`

## Usage

1) Start by declaring a new hook which requires specifying the _type_ of data that it will dispatch as well as a name. This is often done as a global variable:

```go
package user

type User struct {
    ID int
    Name string
    Email string
}

var HookUserInsert = hooks.NewHook[User]("user.insert")
```

2) Dispatch the data to the hook _listeners_:

```go
func (u *User) Insert() {
    db.Insert("INSERT INTO users ...")
    
    HookUserInsert.Dispatch(u)
}
```

Or, dispatch all listeners asynchronously with `HookUserInsert.DispatchAsync(u)`.

3) Listen to a hook:

```go
package greeter

func init() {
    user.HookUserInsert.Listen(func(e hooks.Event[user.User]) {
        sendEmail(e.Msg.Email)
    })
}
```

### Things to know

- The `Listen()` callback does not have to be an anonymous function. You can also do:

```go
package greeter

func init() {
    user.HookUserInsert.Listen(onUserInsert)
}

func onUserInsert(e hooks.Event[user.User]) {
    sendEmail(e.Msg.Email)
}
```

- If you are using `init()` to register your hook listeners and your package isn't being imported elsewhere, you need to import it in order for that to be executed. You can simply include something like `import _ "myapp/greeter"` in your `main` package.
- The `hooks.Event[T]` parameter contains the data that was passed in at `Event.Msg` and the hook at `Event.Hook`. Having the hook available in the listener means you can use a single listener for multiple hooks, ie:

```go
HookOne.Listen(listener)
HookTwo.Listen(listener)

func listener(e hooks.Event[SomeType]) {
    switch e.Hook:
        case HookOne:
        case HookTwo:
}
```

- Since the `Msg` is provided as a _pointer_, a hook can modify the the data which can be useful to allow for modifications prior to saving a user, for example.
- You do not have to use `init()` to listen to hooks. For example, another pattern for this example could be:

```go
package greeter

type Greeter struct {
    emailClient email.Client
}

func NewGreeter(client email.Client) *Greeter {
    g := &Greeter{emailClient: client}
    
    user.HookUserInsert.Listen(func (e hooks.Event[user.User]) {
        g.sendEmail(e.Msg.Email)
    })
    
    return g
}
```

## More examples

While event-driven usage as showed above is the most common use-case of hooks, they can also be used to extend functionality and logic or the process in which components are built. Here are some more examples.

### Router construction

If you're building a web service, it could be useful to separate 

### Modifications

### Validation