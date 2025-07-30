# Junkyard's backend

## Intro

This repo contains the codebase for the backend server of the currently-unnamed social media app.
Written in Golang, each directory is a separate Golang package, each serving a different purpose.
The `main` package, as its name suggests, is the one responsible for launching the server.
Package `config` takes care of reading and loading the server's configuration, which includes reading environment variables like server host, port, the router package used, and so on.
The `api` package takes care of creating (REST-ful) API handlers and serves incoming requests.
`database` is the package that connects to an external storage service (e.g. Postgres).
`util` is a utility package, so it contains functions and structs that can be used in other packages.

## Quickstart

The easiest way to start the server is with [docker](docker.com). If available, the command `docker compose build && docker compose up` ran in a Linux terminal (in this repo's directory) will build and run both the server and the external storage service it will connect to.
Doing this is only recommended for development. Production will require a more robust solution.
Alternatively, it is possible to run just the server with the help of the Makefile: this requires the `make` package for the terminal; if available, then running `make run.dev` or `make run.prod` will start *just the server* in either development or production mode.
Note that in this case, the storage must be started separately (e.g. as a docker container or locally).
