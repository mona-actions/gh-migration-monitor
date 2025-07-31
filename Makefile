# Makefile for gh-migration-monitor

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Binary name
BINARY_NAME=gh-migration-monitor
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe

# Build flags
LDFLAGS=-ldflags="-w -s"

# Test flags
TEST_FLAGS=-v -race -coverprofile=coverage.out

.PHONY: all build clean test coverage lint fmt help install deps

all: test lint build ## Run tests, lint, and build

help: ## Show this help message
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download and verify dependencies
	$(GOMOD) download
	$(GOMOD) verify
	$(GOMOD) tidy

build: ## Build the binary
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v

build-all: ## Build for all platforms
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/linux-amd64/$(BINARY_NAME)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/linux-arm64/$(BINARY_NAME)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/darwin-amd64/$(BINARY_NAME)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/darwin-arm64/$(BINARY_NAME)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/windows-amd64/$(BINARY_WINDOWS)

test: ## Run tests
	$(GOTEST) $(TEST_FLAGS) ./...

test-verbose: ## Run tests with verbose output
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...

coverage: test ## Generate test coverage report
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

bench: ## Run benchmarks
	$(GOTEST) -bench=. -benchmem ./...

lint: ## Run linter
	$(GOLINT) run

lint-fix: ## Run linter and fix issues automatically
	$(GOLINT) run --fix

fmt: ## Format code
	$(GOFMT) -s -w .
	$(GOCMD) mod tidy

fmt-check: ## Check if code is formatted
	@test -z $$($(GOFMT) -l .) || (echo "Code is not formatted. Run 'make fmt'" && exit 1)

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)
	rm -f coverage.out
	rm -f coverage.html
	rm -rf dist/

install: build ## Install the binary to GOPATH/bin
	cp $(BINARY_NAME) $(GOPATH)/bin/

run: build ## Build and run the application
	./$(BINARY_NAME)

dev-setup: ## Set up development environment
	@echo "Setting up development environment..."
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development tools installed successfully!"

check: fmt-check lint test ## Run all checks (format, lint, test)

ci: deps check build ## Run CI pipeline locally

# Docker targets (if needed in future)
docker-build: ## Build Docker image
	docker build -t $(BINARY_NAME) .

docker-run: ## Run Docker container
	docker run --rm -it $(BINARY_NAME)

.DEFAULT_GOAL := help
