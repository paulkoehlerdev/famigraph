# Default target when no arguments are given
.DEFAULT_GOAL := help

# Explicitly declare our targets as PHONY
.PHONY: help build test clean deploy

# Colors for terminal output
BLUE := \033[36m
RESET := \033[0m

# Parse comments after targets and display them in help
help:
	@echo "Available targets:"
	@awk '/^[a-zA-Z0-9_-]+:.*?## .*$$/ { \
		printf "  $(BLUE)%-20s$(RESET) %s\n", \
		substr($$1, 1, length($$1)-1), \
		substr($$0, index($$0, "##") + 3) \
	}' $(MAKEFILE_LIST)

run: build ## Run the application
	@echo "Running the application..."
	@./build/famigraph

dev: ## run in dev mode with air
	@air -c .air.toml

build: ## Build the application
	@echo "Building the application..."
	@go build -o build/famigraph -ldflags='-s -w -X "main.version=dev"' ./cmd

lint: ## Lint the project
	@echo "Linting..."
	@golangci-lint run

test: ## Run all tests
	@echo "Running tests..."
	@go test ./...

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	# TODO implement

deploy: lint test build ## Deploy to production
	@echo "Deploying to production..."
	ansible-playbook ansible/deploy.yml -i ansible/inventory --ask-vault-pass
