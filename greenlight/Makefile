## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
 	
	
# confirm target.
# [ $${ans:-N} = y ] to evaluate response, this will return true if the user enters y and false otherwise
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]
	
## run/api: run the cmd/api application
run/api:
	go run ./cmd/api
	
## db/psql: connect to the database using psql
db/psql:
	psql ${GREENLIGHT_DB_DSN}
	
## db/migrations/new name=$1: create a new database migration
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}
	
# Include the confirm target as prerequisite.
## db/migrations/up: apply all up database migrations
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up
