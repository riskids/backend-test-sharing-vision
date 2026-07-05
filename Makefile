.PHONY: all build run test clean setup-db migrate help

# Variables
APP_NAME=article-microservice
APP_PORT=8080
DB_NAME=article_db
DB_USER=root

# Default target
all: build

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) ./cmd/api
	@echo "Build complete: bin/$(APP_NAME)"

# Run the application
run:
	@echo "Starting $(APP_NAME) on port $(APP_PORT)..."
	@go run ./cmd/api/main.go

# Run tests using the test script
test: build
	@echo "Running API tests..."
	@chmod +x scripts/test_api.sh
	@./scripts/test_api.sh

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@echo "Clean complete"

# Setup database (requires MySQL installed)
setup-db:
	@echo "Setting up database..."
	@echo "Please run the following commands manually:"
	@echo "  mysql -u $(DB_USER) -p -e 'CREATE DATABASE IF NOT EXISTS $(DB_NAME);'"
	@echo "  mysql -u $(DB_USER) -p $(DB_NAME) < migrations/000001_create_posts_table.up.sql"

# Run database migration
migrate:
	@echo "Running database migration..."
	@mysql -u $(DB_USER) -p $(DB_NAME) < migrations/000001_create_posts_table.up.sql
	@echo "Migration complete"

# Development mode with auto-reload (requires air installed)
dev:
	@echo "Running in development mode..."
	@go run ./cmd/api/main.go

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

# Help
help:
	@echo "Available targets:"
	@echo "  make build      - Build the application"
	@echo "  make run        - Run the application"
	@echo "  make test       - Build and run API tests"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make setup-db   - Show database setup instructions"
	@echo "  make migrate    - Run database migration"
	@echo "  make deps       - Install dependencies"
	@echo "  make dev        - Run in development mode"