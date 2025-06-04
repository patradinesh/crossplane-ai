# Crossplane AI Makefile

# Variables
BINARY_NAME=crossplane-ai
VERSION?=0.1.0
BUILD_DIR=build
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

.PHONY: all build clean test coverage deps help install uninstall run

## help: Show this help message
help:
	@echo "Available commands:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## all: Build the project
all: clean build

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

## build-all: Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## test-race: Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	$(GOTEST) -race -v ./...

## coverage: Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## lint: Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

## install: Install the binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	cp $(BINARY_NAME) $(GOPATH)/bin/

## uninstall: Remove the binary from GOPATH/bin
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	rm -f $(GOPATH)/bin/$(BINARY_NAME)

## run: Build and run the application
run: build
	./$(BINARY_NAME)

## demo: Run demo commands
demo: build
	@echo "Running demo commands..."
	@echo "\n🤖 Asking about AWS resources:"
	./$(BINARY_NAME) ask "what AWS resources do I have?"
	@echo "\n💡 Getting database suggestions:"
	./$(BINARY_NAME) suggest database
	@echo "\n🔬 Analyzing resources:"
	./$(BINARY_NAME) analyze --summary
	@echo "\n📝 Generating a database manifest:"
	./$(BINARY_NAME) generate "MySQL database on AWS" --dry-run

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t crossplane-ai:$(VERSION) .

## release: Create a release
release: test build-all
	@echo "Creating release $(VERSION)..."
	@echo "Binaries built in $(BUILD_DIR)/"

## dev: Start development mode
dev: build
	@echo "Starting development mode..."
	./$(BINARY_NAME) --verbose

# Default target
.DEFAULT_GOAL := help
