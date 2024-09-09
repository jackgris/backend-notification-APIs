# Include variables from the .env file
include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## run/db: run database and pgadmin
.PHONY: run/db
run/db:
	docker compose up

## run/api: run the API
.PHONY: run/api
run/api:
	@source ./.env && go run ./cmd/api

## db/migration/up: apply all up database migrations
.PHONY: db/migration/up
db/migration/up: confirm
	@echo 'Running up migrations...'
	migrate -database ${DATABASE_URL} -path ./migrations up

## db/migration/down: apply all down database migrations
.PHONY: db/migration/down
db/migration/down: confirm
	@echo 'Running down migrations...'
	migrate -database ${DATABASE_URL} -path ./migrations down

