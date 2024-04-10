
DB_DSN = "postgres://chat:chat_db_password@127.0.0.1:5432/chat?sslmode=disable"

confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## db/migrations/new name=$1: create a new database migration
db/migrations/new: 
	echo "creating migration files for ${name}.."
	# migrate create -ext=.sql -format="2006-01-02_15-04-05" -seq -dir=./migrations ${name}
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
db/migrations/up:
	echo "Running up migrations..."
	@migrate -path ./migrations -database ${DB_DSN} -verbose up

## db/migrations/:force force fixing the migration version
db/migrations/force: confirm
	@echo "Force fixing the failed migration number: ${force}"
	@migrate -path ./migrations -database ${DB_DSN} force ${force}

db/init:
	./migrations/init.sh

.PHONY: confirm db/migrations/new db/migrations/force db/init
