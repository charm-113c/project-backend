build:
	@echo "Compiling packages"
	@go build -o bin/backend ./main

run.prod: build
	@echo "Running in prod mode"
	@DEV_MODE=false ./bin/backend

run.dev: build
	@echo "Running in dev mode"
	@DEV_MODE=true ./bin/backend

test:
	@go test -v ./...
