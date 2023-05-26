## help: print this help message

.PHONY: help
help: 
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
	
# confirm target.
# [ $${ans:-N} = y ] to evaluate response, this will return true if the user enters y and false otherwise
#
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]
	
## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api
	
## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${GREENLIGHT_DB_DSN}
	
## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}
	
# Include the confirm target as prerequisite.
## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up


# A phony target is one that is not really the name of a file; rather it is just a name for a rule 
# to be executed. so i.g if you have the run/api command above, and you have for example a folder that has 
# the same path (i.g  run/api/main.go.... ) than make will not success in executing this rule, and will see
# something like  (make: 'run/api' is up to date.) so phony prevent this behavior 

