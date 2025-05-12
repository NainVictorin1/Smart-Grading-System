## Filename Makefile
include .envrc

# Testing - Run the tests

PHONY: run migrateup resetdb


migrateup:
	@echo "Applying migrations..."
	goose -dir $(MIGRATION_PATH) postgres "$(DB_DSN)" up

resetdb:
	@echo "Resetting database..."
	-dropdb -U $(DB_USER) --if-exists $(DB_NAME)
	createdb -U $(DB_USER) $(DB_NAME)
	$(MAKE) migrateup


run:
	@echo "Starting application..."
	CSRF_KEY=$(CSRF_KEY) \
	SESSION_SECRET=$(SESSION_SECRET) \
	DB_DSN=$(DB_DSN) \
	go run cmd/web/main.go

test:
	@echo "Running tests..."
	go test -v -cover ./...

lint:
	@echo "Linting code..."
	golangci-lint run

clean:
	@echo "Cleaning up..."
	go clean
	rm -f coverage.out

# Database management shortcuts
createdb:
	createdb $(DB_NAME)

dropdb:
	dropdb $(DB_NAME)

migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir $(MIGRATION_PATH) create $${name} sql
# Helpers
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out