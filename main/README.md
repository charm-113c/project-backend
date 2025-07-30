# The main package

Being the main, it's responsible for starting up the server and ensuring it's listening to and serving requests.

## Loading a configuration

The default configuration is found in `config` package/directory, more specifically in `config.go` in the `Config` struct. A configuration can be set through the configuration file in `../config/config.yaml`, as well as through environmental variables (through the `../config/server.env` file, or alternatively through the CLI). Note: the config files are .gitignored for security, but the struct allows to easily reconstruct the structure of YAML config files. Alternatively, config files can be omitted, thus leaving the default values in the `Config` struct.  
Configurations are read with the [Clean Env](https://github.com/ilyakaznacheev/cleanenv) package. This package first reads the values from the config files and parses them into the `Config` struct, then reads environmental variables and overwrites the values from the file.  
The keys to be inputted can be found in the same `Config` struct found in `config.go`; the same struct also indicates the structure for YAML config files.

Note: one of the fields is called `DevMode`. This field is a boolean that specifies whether the server is to be run in development mode (in which case its value is `true`) or in production mode (in which case it's `false`).  
Development mode starts a development logger and sets the log file in this project's root directory so that it can be easily read, while production mode starts a production logger and writes the log file in `/var/log/`, a more appropriate directory.  
This field does *not* need to be set in the config files, more on that in the [Starting the server](##starting-the-server) section.

## Starting the server

### Without Docker

The server can be easily started with the help of the Makefile. The Makefile defines two commands to start the server either in production or development mode, without having to input environmental variables. `make run.prod` will start the server in production mode, while `make run.dev` will set the `DEVMODE` env variable to true and then start the server in development mode.  
Tests can also be ran with `make test`, which makes use of the Golang [testing](https://pkg.go.dev/testing) package.
Note: doing this only starts the server and **not** the storage/database. A separate Docker container could be used to initialise the storage, or the solution below can be followed.

### With Docker and Docker Compose

For development's sake, a `dev.Dockerfile` as well as a `docker-compose.yaml` were written.
This allows standardising the development environment on top of starting the server and the database in a single command.
The command in question is: `docker compose build && docker compose up`.
Note: for the time being, the `docker-compose` is geared towards Postgres. To run other storage solutions, modifications
are required.

## Code flow

As the code structure itself is subject to change, only the general structure and flow is reported here:

- The `main()` function calls the auxiliary `run()` function to start the server
- The configuration is then loaded and validated
- A logger is started and its closing is deferred; using a dedicated logger allows writing to one or more files rather than just to the console
- A connection is established with the storage resources
- A router is initialized
- The router starts listening to and serves HTTP requests
- At the same time, `run()` listens to external shutdown signals (e.g. interrupts a.k.a. CTRL+C)
- If server errors or shutdown signals are received, a graceful shutdown is initiated
- The connection to the storage objects are closed
- The program exits
