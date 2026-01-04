.PHONY: help build-dev build-prod build-darwin build-linux build-windows clean ent-generate ent-diff ent-apply ent-status

# Variables
APP_NAME=crud-solution
DEV_BINARY=bin/$(APP_NAME)-dev
PROD_BINARY=bin/$(APP_NAME)

# Database connection for local development
DB_USER=crud
DB_PASS=crud@2026
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=crud-db
DEV_URL=mysql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?parseTime=true

help:
	@echo "Available targets:"
	@echo "  build-dev          - Build development binary (dev tag)"
	@echo "  build-prod         - Build production binary (prod tag)"
	@echo "  build-darwin       - Build for macOS (darwin/arm64)"
	@echo "  build-linux        - Build for Linux (linux/amd64)"
	@echo "  build-windows      - Build for Windows (windows/amd64)"
	@echo "  clean              - Remove build artifacts"
	@echo "  ent-generate       - Generate Ent code from schema"
	@echo "  ent-diff           - Generate migration files"
	@echo "  ent-apply          - Apply migrations to database"
	@echo "  ent-status         - Check migration status"

# Development build (local platform)
build-dev:
	@echo "Building development binary..."
	@mkdir -p bin
	go build -tags dev -o $(DEV_BINARY) .

# Production build (local platform)
build-prod:
	@echo "Building production binary..."
	@mkdir -p bin
	go build -tags prod -o $(PROD_BINARY) .

# Cross-platform builds
build-darwin:
	@echo "Building for macOS (darwin/arm64)..."
	@mkdir -p bin
	GOOS=darwin GOARCH=arm64 go build -tags prod -o bin/$(APP_NAME)-darwin-arm64 .

build-linux:
	@echo "Building for Linux (linux/amd64)..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -tags prod -o bin/$(APP_NAME)-linux-amd64 .

build-windows:
	@echo "Building for Windows (windows/amd64)..."
	@mkdir -p bin
	GOOS=windows GOARCH=amd64 go build -tags prod -o bin/$(APP_NAME)-windows-amd64.exe .

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf tmp/
	@rm -rf crud_temp/

# Ent operations (require local Atlas installation)
ent-generate:
	@echo "Generating Ent code..."
	go generate ./internal/database/ent

ent-diff:
	@echo "Generating migration diff..."
	atlas migrate diff --dir file://internal/database/migrations --to ent://internal/database/ent/schema --dev-url "$(DEV_URL)"

ent-apply:
	@echo "Applying migrations..."
	atlas migrate apply --dir file://internal/database/migrations --url "$(DEV_URL)"

ent-status:
	@echo "Checking migration status..."
	atlas migrate status --dir file://internal/database/migrations --url "$(DEV_URL)"
