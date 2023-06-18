## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## run: run the cmd application
.PHONY: run
run:
	cd public && ./tailwindcss -i styles.css -o global.css --minify
# ./tailwindcss -i styles.css -o global.css --watch
	CGO_ENABLED=0 go run . -db-dsn="db/db.db"

## db/sqlite: connect to the database using sqlite
.PHONY: db/sqlite
db/sqlite:
	sqlite3 db/db.db

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./db/migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./db/migrations -database sqlite://db/db.db up

## audit: tidy dependencies, format, vet and test all code
.PHONY: audit
audit:
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify