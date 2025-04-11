## Filename Makefile
include .envrc

# Testing - Run the tests
.PHONY: run/tests
run/tests: vet
	go test -v ./...

# Format - Format the code
.PHONY: fmt
fmt: 
	go fmt ./...

# Vet - Run static analysis on the code
.PHONY: vet
vet: fmt
	go vet ./...

# Run the application
.PHONY: run
run: vet
	go run ./cmd/web -addr=${ADDRESS} -dsn=${GRADE_DB_DSN}

# Start PostgreSQL shell
.PHONY: db/psql
db/psql:
	psql ${GRADE_DB_DSN}

# Create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

# Apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GRADE_DB_DSN} up

# Undo the last migration
.PHONY: db/migrations/down-1
db/migrations/down-1:
	@echo 'Running down migrations...'
	migrate -path ./migrations -database ${GRADE_DB_DSN} down 1

# Fix a dirty SQL migration
.PHONY: db/migrations/fix
db/migrations/fix:
	@echo 'Checking migration status...'
	@migrate -path ./migrations -database ${GRADE_DB_DSN} version > /tmp/migrate_version 2>&1
	@cat /tmp/migrate_version
	@if grep -q "dirty" /tmp/migrate_version; then \
		version=$$(grep -o '[0-9]\+' /tmp/migrate_version | head -1); \
		echo "Found dirty migration at version $$version"; \
		echo "Forcing version $$version..."; \
		migrate -path ./migrations -database ${GRADE_DB_DSN} force $$version; \
		echo "Running down migration..."; \
		migrate -path ./migrations -database ${GRADE_DB_DSN} down 1; \
		echo "Running up migration..."; \
		migrate -path ./migrations -database ${GRADE_DB_DSN} up; \
	else \
		echo "No dirty migration found"; \
	fi
	@rm -f /tmp/migrate_version