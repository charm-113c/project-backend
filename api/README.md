# The API package

IMPORTANT: this README is obsolete, as the router isn't abstracted anymore and the Echo package is used directly. While updates have been made in the parent repo, they have not (and likely will not) made their way here.

The api package is responsible for defining and implementing the API requests a client can make. It is designed to be modular and flexible by decoupling the functions and methods that need to be implemented from their actual implementation. Put simply, it defines structs and interfaces that contain the various methods that HTTP routers will need to implement. As such, any router package can be chosen for the implementation (e.g. Fiber, Echo, Gin, etc...), and the router to be activated when the server is started will be decided through the configuration files and/or environmental variables.

## The api.go file

This is the main file. It defines the `Router` struct, the main abstraction between the HTTP router the server will instantiate and its implementation. The `StartRouter` function, as its name indicates, is responsible for instantiating the `Router` that will be used, based on the provided configuration.

### The Router struct

As mentioned above, this is the main layer of abstraction. This struct determines what any implementation of a Router needs to do.
It has three components: CoreRouter, Logger and Handle.

- Logger: The simplest, a logger dedicated to the router.
- CoreRouter: an interface that defines fundamental router methods like Listen() or ListenTLS(), which for instance make the router listen to HTTP requests on a given address.
- Handle: an inner struct, it lists out all the requests that will need to be handled. For example, Router.Handle.UserReqs.GetProfile() would handle a request to get a user's profile.

#### More details on the Router.Handle field

One of the fields of the `Router` struct is `Handle`, of type `RequestHandlers`. `RequestHandlers` is itself a struct, which contains several interfaces: these interfaces define the methods that the router will need to handle. For example, the `UserRequests` interface defines the method `createUser()`: any implementation of `Router` will therefore need to have a `Handle` field that implements this method. In short, we must be able to call `Router.Handle.UserRequests.createUser()`, a rather self-explanatory name.
These interfaces are not only important, they can also be very large, as the number and type of possible requests can grow over time. As such, each one of these interfaces has its own file, with the name `*intrfc-...-requests.go*`.
For the sake of clarity, it should be mentioned that the interfaces' methods are created and named after the operationID groups in the documentation repository (docs/operationIDs.md), and their methods correspond to the different operationIDs within each group.

### The StartRouter function

Called by the `main` package, this function creates the Router instance to be used based on an empty `Router` struct given as input.

#### Code Flow (subject to change)

**Code flow out of date, abstraction removed in favour of simplicity (2025-07-10)**

- Given parent logger, create child logger specific to the Router struct
- Maximise number of connections that this instance of the server can handle: Linux systems are usually limited to some 1000 file descriptors, but each connection requires a file descriptor, so we maximise the number of file descriptors that can be used by the system.
- Implement the `Router.CoreRouter` according to the configuration
- Set up the HTTP routes (e.g. start listening to and serve GET requests on route /users/{userID}/profile). [!NOTE] Since there are a lot of routes, the actual function listing them is put in a file of its own

## The components.go file

This file simply defines the various structs that are needed to receive and serve requests. As a simple example, the `CreateUser()` method will take as input a `NecessaryUserData` struct defined in this file; as its name implies, this struct just contains the fields necessary to create an account. `NecessaryUserData` actually follows the structure defined in the documentation repository (docs/OpenAPI/Account-Profile-APIs.yaml) for the createUser operationID. In other words, it follows the definition of the request body.
In essence, this file contains the structs of the request and response bodies to the API requests.

## From abstraction to implementation

This begins in `StartRouter`, however it's worth explaining what's happening in more details.
`StartRouter` is given an empty `Router` object as an input, and is responsible for populating it.

### Implementing `Router.Logger`

The first field is the logger, inherited from a parent logger that is also given in input to the function.

### Implementing `Router.CoreRouter`

The second field is `CoreRouter`. The router package used will be put in this field (e.g. Fiber, Echo, Gin) when the function `New<package>Router` is called.
The latter takes as input the `Router` object, and replaces its currently empty `CoreRouter` field with a struct implementing the `CoreRouterInterface` interface, effectively giving us a valid `CoreRouter`.
In general, that struct will be a wrapper over the package's main router in order to endow it with new methods while retaining the package's methods.

### Implementing `Router.Handle`

This is the most involving step. Using the methods of `CoreRouter`, we now implement all the routes of the API requests we must serve. As a simple example, using `CoreRouter`'s `AddRoute` method, we can add a route to handle the request of an account creation by calling `CoreRouter.AddRoute("POST", "/users/", createUser)`.
To define which routes must be implemented, we do that indirectly by defining the handler functions that should handle them (in the above example, `createUser`).
In order to define those in an organised manner, the `Router.Handle` field was made into a struct. However, it's a special struct in that all its fields are interfaces.
Those interfaces divide the handlers in categories, e.g. `UserRequests` is the interface containing handlers relating to user and user account related requests, `EventRequest` does the same for events, etc.
As such, to implement `Router.Handle`, the corresponding struct must be implemented, and to do that all of its interfaces must be implemented.

TODO: use this to explain

```go
// api/handlers.go
type RequestHandlers struct {
    User   *UserHandlers
    Social *SocialHandlers
    // ... other domains
}

type UserHandlers struct {
    DB     database.UserStorageHandler
    Logger *zap.Logger
}

func NewUserHandlers(db database.UserStorageHandler, logger *zap.Logger) *UserHandlers {
    return &UserHandlers{DB: db, Logger: logger}
}

// Example endpoint method
func (h *UserHandlers) CreateUser(c echo.Context) error {
    // Business logic using h.DB and h.Logger
    return c.JSON(200, "User created")
}
```
