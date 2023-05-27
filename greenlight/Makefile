include .envrc

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

# ==================================================================================== #
# 		DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@go run ./cmd/api -db-dsn=${GREENLIGHT_DB_DSN} -jwt-secret=${JWT_SECRET}

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${GREENLIGHT_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up

# ==================================================================================== #
# 		QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...
	
## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor


# ==================================================================================== #
#  		BUILD
# ==================================================================================== #
#
current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	# -ldflags="-s"   to reduce the binary size
	# using the -X linker flag to ‘burn in’ the current time to our main.buildTime variable 
	@echo 'Building cmd/api...'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api


# ==================================================================================== #
# 		PRODUCTION
# ==================================================================================== #
production_host_ip = '45.55.49.87'

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh greenlight@${production_host_ip}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	rsync -P ./bin/linux_amd64/api greenlight@${production_host_ip}:~
	rsync -rP --delete ./migrations greenlight@${production_host_ip}:~
	rsync -P ./remote/production/api.service greenlight@${production_host_ip}:~
	rsync -P ./remote/production/Caddyfile greenlight@${production_host_ip}:~
	ssh -t greenlight@${production_host_ip} '\
	migrate -path ~/migrations -database $$GREENLIGHT_DB_DSN up \
	&& sudo mv ~/api.service /etc/systemd/system/ \
	&& sudo systemctl enable api \
	&& sudo systemctl restart api \
	&& sudo mv ~/Caddyfile /etc/caddy/ \
	&& sudo systemctl reload caddy \
	'

# -t flag to run these database migrations on our server as the greenlight user
# we used $$ to escape the first dollar,,,, It’s also important to note that we’re surrounding this 
# command with single quotes. If we used double quotes, it would be an interpreted string and we 
# would need to use an additional escape character \ like so:
# "migrate -path ~/migrations -database \$$GREENLIGHT_DB_DSN up"


# A phony target is one that is not really the name of a file; rather it is just a name for a rule 
# to be executed. so i.g if you have the run/api command above, and you have for example a folder that has 
# the same path (i.g  run/api/main.go.... ) than make will not success in executing this rule, and will see
# something like  (make: 'run/api' is up to date.) so phony prevent this behavior 

