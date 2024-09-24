# Project name
PROJECT_NAME := go-tempalte

# Go-related variables
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)

# Swag-related variables
SWAG := swag

# Main Go package and file
MAIN_PACKAGE := ./cmd/api
MAIN_FILE := $(MAIN_PACKAGE)/main.go

# Golang Migrate
MIGRATE := migrate

# Read database connection details from config.yaml
DB_HOST := $(shell yq e '.database.host' config.yaml)
DB_PORT := $(shell yq e '.database.port' config.yaml)
DB_USER := $(shell yq e '.database.username' config.yaml)
DB_PASS := $(shell yq e '.database.password' config.yaml)
DB_NAME := $(shell yq e '.database.name' config.yaml)

# Construct DB_URL
DB_URL := "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

# Migration files directory
MIGRATIONS_DIR := ./migrations

.PHONY: all build run clean docs help migrate-create migrate-up migrate-down migrate-force

all: build

# Build the project
build:
	@echo "Building $(PROJECT_NAME)..."
	@go build -o $(GOBIN)/$(PROJECT_NAME) $(MAIN_FILE)

# Run the project
run:
	@go run $(MAIN_FILE)

# Clean build files
clean:
	@echo "Cleaning..."
	@rm -rf $(GOBIN)
	@rm -rf docs

# Install swag
install-swag:
	@echo "Installing swag..."
	@go install github.com/swaggo/swag/cmd/swag@latest

# Generate API documentation
docs: install-swag
	@echo "Generating API documentation..."
	@$(SWAG) init --parseDependency --parseInternal -g $(MAIN_FILE) -o ./docs

# Run tests
test:
	@go test ./...

migrate-create:
	@read -p "Enter migration name: " name; \
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $${name}

migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database $(DB_URL) up

migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database $(DB_URL) down 1

migrate-force:
	@read -p "Enter version to force: " version; \
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database $(DB_URL) force $$version

# Display help information
help:
	@echo "Available commands:"
	@echo "  make build    - Build the project"
	@echo "  make run      - Run the project"
	@echo "  make clean    - Clean build files"
	@echo "  make docs     - Generate API documentation using swag"
	@echo "  make test     - Run tests"
	@echo "  make help     - Show this help message"
	@echo "  make migrate-create - Create a new migration"
	@echo "  make migrate-up     - Run all migrations"
	@echo "  make migrate-down   - Rollback the last migration"
	@echo "  make migrate-force  - Force a migration by version"