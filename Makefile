include .envrc

.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## run the server
.PHONY: run/api
run/api:
	go run ./cmd/main -db-dsn=${GREENLIGHT_DB_DSN}

current_time = $(shell date --iso-8601=seconds)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=$(shell git describe --tags --always --dirty --long)'
# build/api: Build the application binary, strip debug information, symbol table
# to reduce binary size. Also build a linux binary.
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags=${linker_flags} -o=./bin/main ./cmd/main
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/main ./cmd/main

## connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${GREENLIGHT_DB_DSN}

## create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up

# apply all down database migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Running down migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} down

# audit: Tidy and verify module dependencies, format and vet code, and run tests.
.PHONY: audit
audit: vendor
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

# vendor: Tidy and verify module dependencies, and vendor them.
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor