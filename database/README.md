# The database package
It takes care of all communications to an external database and (possibly internal) cache. Like the other packages, this one is also made to be flexible and modular; it defines interfaces that contain the methods that will be needed to serve requests (e.g. GetUserData), but it provides a layer of abstraction over the implementation: any solution that implements the interfaces can be used, be it MongoDB, Postgres, MySQL, etc...

## database.go
The `database/database.go` file is the main file of this package. It is designed to provide a layer of abstraction between the queries to the database and the actual implementation of the database and its queries. Put simply, `database/database.go` defines the operations in interfaces, so that `main/main.go` only needs to worry about calling those interfaces, not how they are actually implemented. This allows for flexibility in the database choice.  

## The storage
The database package defines a `Storage` struct, made up of multiple components:

- an inner struct called `Conns`
- a DB-specific logger
- a cache

This struct is the abstraction between the database queries and the implementation of these queries.

### `Storage.Conns`
`Conns` stands for connection pool: in order to maximize a DB's performance, it is better to establish multiple connections with it rather than use a single connection. As such, any implementation of the DB should allow for a connection pool. `Conns` is a struct that contains multiple interfaces: these interfaces define the operations that the DB must implement, they are the queries that the backend needs to make to the database in order to respond to requests.  
There are multiple interfaces in order to keep things more organized: for example, the `UserStorageHandler` interface represents the operations that can be done to the User table, e.g.: GetUserByID, CreateUser, UpdateUserInfo, etc...

### The logger
This logger is a child of the logger from `main/main.go`, and as such inherits its configuration. It has an added field, `{component: database}` and is otherwise independent from its parent logger.

### The cache
Golang routers are highly performant, with the majority of them capable of handling thousands of requests per second and some even tens of thousands. However, in our case the routers must make database queries, which can be and often are slow. This would nullify the performance of Golang routers. As such, having a cache in front of the database would decrease interactions with the database and thus reduce its impact on performance.

## Implemented DBs
### Postgres
Being one of the most mature and popular DBs, Postgres is the first choice. Performant, scalable vertically -and with extensions, horizontally- it is *the* general-purpose SQL database, and as such is a shoo-in for this project. Until our needs are clarified and a more suitable tool is found, there isn't a reason to not choose Postgres.
