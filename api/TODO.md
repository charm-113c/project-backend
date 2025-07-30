# TODO

## 16/05/25

The router abstraction will be eliminated, Echo package will be used directly.
We'll need to re-organise the code to do this: initialising the router can be
done in the `api` package, but starting the HTTP/S server can be done directly
in `main/main.go`.

In initialising the Echo router, we'll create subhandlers to be responsible
for routing groups, e.g. EventHandler for EventRequests.

``` go
// api/handlers.go
type RequestHandlers struct {
    // The fields themselves will be interfaces, as we've already defined
    // them in api.go, however their implementations will be as seen here
    User   *UserHandlers
    Social *SocialHandlers
    // ... other domains
}

type UserHandlers struct {
    DB     database.DB
    Logger *zap.Logger
}

func NewUserHandlers(db database.DB, logger *zap.Logger) *UserHandlers {
    return &UserHandlers{DB: db, Logger: logger}
}

// Example endpoint method
func (h *UserHandlers) CreateUser(c echo.Context) error {
    // Business logic using h.DB and h.Logger
    return c.JSON(200, "User created")
}

//------

// main.go
func setupHandlers(cfg config.Config, db *database.Storage, logger *zap.Logger) *api.RequestHandlers {
    return &api.RequestHandlers{
        User:   api.NewUserHandlers(db, logger),
        Social: api.NewSocialHandlers(db, logger),
        // ...
    }
}
```

We will also need to organise the dependency injection of the DB.
In the above code it looks like we give the full DB to each subhandler,
but in reality it is better to separate concerns here too; i.e. each
subhandler only has access to what it needs:

``` go
// database.go
type Storage struct {
    userConn   UserStorageHandler
    socialConn SocialStoragesHandler
    // ...
}

func (s *Storage) User() UserStorageHandler {
    return s.userConn
}

func (s *Storage) Social() SocialStoragesHandler {
    return s.socialConn
}

//--------

// handlers/user.go
func (h *UserHandlers) CreateUser(c echo.Context) error {
    // Access ONLY the user DB
    if err := h.DB.User().CreateUser(...); err != nil {
        // Handle error
    }
}
```

## 17/05/2025

We need to re-organise the DB and how it's handled, in order to accommodate for the above dependency injection.

## 18/05/2025

Turns out the DB is already primed for that. However, we need to handle caching too. So, we need to set up the ground to make use of the caching interface.
We'll use a decorator pattern over each of the storage handlers so that the router subhandlers remain blissfully unaware of whether they're
getting their data from the cache or from the DB.

## 21/05/2025

Since subhandlers shouldn't know whether they're searching the cache or not, we'll first deal with them and how they'll respond to requests.

## 2025-07-10

Create middleware for authentication.
The case of account creation being potentially special, we need to understand what's necessary. We need a refresher on OAuth2.0 too.

In the future, we'll also want *an analytics middleware*

## 2025-07-11

Authentication is less of a pain than we thought, but we have to be careful about authorisation.
We need to find a way to bind access tokens to a device. There's RFC 8471 for that, but that might be overkill.
If you can think of a simple but efficient way to do it, we'll do it. However, the idea will likely need to involve
having the client **securely** generate some secret.

## 2025-07-19

Implement logto auth2: create dev-mode local session storage that implements Logto's storage interface, init logto config,
set up auth routes.

## 2025-07-20

Implement storage interface, implement the sign-in, callback, and potentially sign-out routes

## 2025-07-22

Implement the logto routes

## 2025-07-24

Test if this works.

## 2025-07-26

It could be the case that to test, we have to implement Logto's sign-in and callback functions. That shouldn't be necessary though.
Figure out why after starting the docker compose, we can't connect to the address given to us by Echo

## 2025-07-28

Container shutdown is too fast, server doesn't have time to shutdown gracefully. Find a way around this.
On top of that, finally implement the sign-in and call back handlers, the server is working
