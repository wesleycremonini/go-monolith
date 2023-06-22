include .envrc

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
	ADDR=":3000" DB_DSN=${DB_DSN} REDIS_HOST=${REDIS_HOST} REDIS_PASS=${REDIS_PASS} CGO_ENABLED=0 go run .

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./db/migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./db/migrations -database ${DB_DSN} up

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

## run/minikube: run the project on minikube kubernetes
.PHONY: run/minikube
run/minikube:
	@if ! command -v minikube > /dev/null; then \
		echo "Minikube is not installed. Please install Minikube."; \
		exit 1; \
	fi
	@echo 'Building go-app...'
	DOCKER_BUILDKIT=1 docker build -t go-app . -f docker/Dockerfile.prod --no-cache
	@echo 'Loading go-app...'
	minikube image load go-app
	@echo 'Deploying go-app...'
	minikube kubectl -- apply -f docker/minikube/deployment.yaml
	@echo 'Waiting for pod to be ready...'
	minikube kubectl -- wait pods --all --for condition=Ready --timeout=90s
	@echo 'Forwarding go-app...'
	minikube kubectl -- port-forward service/go 3000:80