# Makefile for Go game project

# Variables
BINARY_NAME=game
GO=go
MAIN_PATH=./game
BUILD_DIR=build

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run without building
.PHONY: run-direct
run-direct:
	@echo "Running directly with go run..."
	$(GO) run $(MAIN_PATH)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Vet code
.PHONY: vet
vet:
	@echo "Vetting code..."
	$(GO) vet ./...

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GO) mod tidy

# Help command
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make build      - Build the application"
	@echo "  make run        - Build and run the application"
	@echo "  make run-direct - Run without building (using 'go run')"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make test       - Run tests"
	@echo "  make fmt        - Format code"
	@echo "  make vet        - Vet code"
	@echo "  make deps       - Install dependencies"
	@echo "  make help       - Show this help message"
