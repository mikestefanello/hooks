# Hooks

## Overview

_Hooks_ provides a simple, **type-safe** hook system to enable easier modularization of your Go code. A _hook_ allows various parts of your codebase to tap in to events and operations happening elsewhere which prevents direct coupling between the producer and the consumers/listeners. For example, a _user_ package/module in your code may dispatch a _hook_ when a user is created, allowing your _notification_ package to send the user an email, and a _history_ package to record the activity without the _user_ module having to call these components directly. A hook can also be used to allow other modules to alter and extend data before it is processed.

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
    Password string
}

var HookUserInsert = hooks.NewHook[User]("user.insert")
```

2) Dispatch the data to the hook _listeners_:

```go
func (u *User) Insert() {
    db.Insert("INSERT INTO users ...")
    
    HookUserInsert.Dispatch(&u)
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
    switch e.Hook {
    case HookOne:
    case HookTwo:
    }
}
```

- If the `Msg` is provided as a _pointer_, a hook can modify the the data which can be useful to allow for modifications prior to saving a user, for example.
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

While event-driven usage as shown above is the most common use-case of hooks, they can also be used to extend functionality and logic or the process in which components are built. Here are some more examples.

### Router construction

If you're building a web service, it could be useful to separate the registration of each of your module's endpoints. Using [Echo](https://github.com/labstack/echo) as an example:

```go
package main

import (
    "github.com/labstack/echo/v4"
    "github.com/myapp/router"
    
    // Modules
    _ "github.com/myapp/modules/todo"
    _ "github.com/myapp/modules/user"
)

func main() {
    e := echo.New()
    router.BuildRouter(e)
    e.Start("localhost:9000")
}
```

```go
package router

import (
    "net/http"
    
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/mikestefanello/hooks"
)

var HookBuildRouter = hooks.NewHook[echo.Echo]("router.build")

func BuildRouter(e *echo.Echo) {
    e.Use(
        middleware.RequestID(),
        middleware.Logger(),
    )
    
    e.GET("/", func(ctx echo.Context) error {
        return ctx.String(http.StatusOK, "hello world")
    })
    
    // Allow all modules to build on the router
    HookBuildRouter.Dispatch(e)
}
```

```go
package todo

import (
    "github.com/labstack/echo/v4"
    "github.com/mikestefanello/hooks"
    "github.com/myapp/router"
)

func init() {
    router.HookBuildRouter.Listen(func(e hooks.Event[echo.Echo]) {
        e.Msg.GET("/todo", todoHandler.Index)
        e.Msg.GET("/todo/:todo", todoHandler.Get)
        e.Msg.POST("/todo", todoHandler.Post)
    })
}
```

### Modifications

Hook listeners can be used to make modifications to data prior to some operation being executed if the _message_ is provided as a pointer. For example, using the `User` from above:

```go
var HookUserPreInsert = hooks.NewHook[*User]("user.pre_insert")

func (u *User) Insert() {
    // Let other modules make any required changes prior to inserting
    HookUserPreInsert.Dispatch(u)
	
    db.Insert("INSERT INTO users ...")
    
    // Notify other modules of the inserted user
    HookUserInsert.Dispatch(u)
}
```

```go
HookUserPreInsert.Listen(func(e hooks.Event[user.User]) {
    // Change the user's name
    e.Msg.Name = fmt.Sprintf("%s-changed", e.Msg.Name)
})
```

### Validation

Hook listeners can also provide validation or other similar input on data that is being acted on. For example, using the `User` again.

```go
type UserValidation struct {
    User User
    Errors []error
}

var HookUserValidate = hooks.NewHook[UserValidation]("user.validate")

func (u *User) Validate() []error {
    uv := &UserValidation{
        User: &u,
        Errors: make([]error, 0)
    }

    if u.Email == "" {
        uv.Errors = append(uv.Errors, errors.New("missing email"))
    }
    
    // Let other modules validate
    HookUserValidate.Dispatch(u)
	
    return uv.Errors
}
```

```go
HookUserValidate.Listen(func(e hooks.Event[user.UserValidate]) {
    if len(e.Msg.User.Password) < 10 {
        e.Msg.Errors = append(e.Msg.Errors, errors.New("password too short"))
    }
})
```

## Logging

By default, nothing will be logged, but you have the option to specify a _logger_ in order to have insight into what is happening within the hooks. Pass a function in to `SetLogger()`, for example:

```go
hooks.SetLogger(func(format string, args ...any) {
    log.Printf(format, args...)
})
```

```
2022/09/07 13:42:19 hook created: user.update
2022/09/07 13:42:19 registered listener with hook: user.update
2022/09/07 13:42:19 registered listener with hook: user.update
2022/09/07 13:42:19 registered listener with hook: user.update
2022/09/07 13:42:19 dispatching hook user.update to 3 listeners (async: false)
2022/09/07 13:42:19 dispatch to hook user.update complete
```