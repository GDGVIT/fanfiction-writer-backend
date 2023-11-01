## help: Display this help message
.PHONY: help
help:
	@echo 'Usage: '
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## run/api: Runs cmd/api
.PHONY: run/api
run/api: 
	@go run ./cmd/api --db-dsn=${FFWRITER_DB_DSN}


## db/psql: connect to database using psql
.PHONY: db/psql
db/psql:
	psql ${FFWRITER_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new: 
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: run all up migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations'
	migrate -path=./migrations -database=${FFWRITER_DB_DSN} up

## db/migrations/down: run all down migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Running all down migrations'
	migrate -path=./migrations -database=${FFWRITER_DB_DSN} down
