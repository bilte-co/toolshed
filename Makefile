# Variables
BINARY_NAME=toolshed
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"
VERSION ?= $(shell git describe --tags --abbrev=0)
NEXT_PATCH := $(shell echo $(VERSION) | awk -F. '{printf "v%d.%d.%d", $$1, $$2, $$3+1}')
NEXT_MINOR := $(shell echo $(VERSION) | awk -F. '{printf "v%d.%d.0", $$1, $$2+1}')
NEXT_MAJOR := $(shell echo $(VERSION) | awk -F. '{printf "v%d.0.0", $$1+1}')

# Build settings
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

# Directories
BUILD_DIR=build
DIST_DIR=dist

.PHONY: help build clean test install uninstall lint fmt vet deps tidy release dev

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: tag
tag:
	@read -p "Enter version tag (e.g. v1.4.0): " VERSION; \
	git tag -a $$VERSION -m "Release $$VERSION"; \
	git push origin $$VERSION; \
	echo "Published $$VERSION"

.PHONY: bump-patch
bump-patch:
	git tag -a $(NEXT_PATCH) -m "Release $(NEXT_PATCH)"
	git push origin $(NEXT_PATCH)

.PHONY: bump-minor
bump-minor:
	git tag -a $(NEXT_MINOR) -m "Release $(NEXT_MINOR)"
	git push origin $(NEXT_MINOR)

.PHONY: bump-major
bump-major:
	git tag -a $(NEXT_MAJOR) -m "Release $(NEXT_MAJOR)"
	git push origin $(NEXT_MAJOR)

build: deps ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@go clean

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

install: build ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) .
	@echo "Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

uninstall: ## Uninstall the binary
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(shell go env GOPATH)/bin/$(BINARY_NAME)

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1; }
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy

dev: deps fmt vet ## Development build and checks
	@$(MAKE) build
	@echo "Development build complete"

release: clean deps fmt vet test ## Build release binaries for multiple platforms
	@echo "Building release binaries..."
	@goreleaser release --clean

# Example usage targets
examples: build ## Show example commands
	@echo "Example commands:"
	@echo ""
	@echo "  # Hash a string"
	@echo "  ./$(BUILD_DIR)/$(BINARY_NAME) hash string 'Hello, World!' --algo sha256"
	@echo ""
	@echo "  # Hash a file"
	@echo "  ./$(BUILD_DIR)/$(BINARY_NAME) hash file README.md --algo sha512"
	@echo ""
	@echo "  # Hash from stdin"
	@echo "  echo 'Hello' | ./$(BUILD_DIR)/$(BINARY_NAME) hash file -"
	@echo ""
	@echo "  # Hash a directory"
	@echo "  ./$(BUILD_DIR)/$(BINARY_NAME) hash dir . --recursive"
	@echo ""
	@echo "  # Compute HMAC"
	@echo "  ./$(BUILD_DIR)/$(BINARY_NAME) hash hmac 'data' --key 'secret'"
	@echo ""
	@echo "  # Validate file checksum"
	@echo "  ./$(BUILD_DIR)/$(BINARY_NAME) hash validate README.md --expected <hash>"
	@echo ""
	@echo "  # Compare two hashes"
	@echo "  ./$(BUILD_DIR)/$(BINARY_NAME) hash compare <hash1> <hash2>"
	@echo ""
	@echo "  # Enable verbose output"
	@echo "  ./$(BUILD_DIR)/$(BINARY_NAME) --verbose hash string 'test'"
