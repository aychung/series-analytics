# ====================================================================================
# Variables
# ====================================================================================
BINARY_NAME=crawler
BUILD_DIR=dist
MAIN_PACKAGE=./cmd/crawler

# Git information for dynamic binary versioning
VERSION=$(shell jj log -r "latest(ancestors(@) & tags())" --template 'tags' --no-graph | grep . || echo "v0.0.0")
COMMIT=$(shell jj log -r "@" --template "commit_id.short(12)" --no-graph 2>/dev/null | grep . || echo "unknown")

# Go compiler configuration flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -s -w"

# ====================================================================================
# Phony Targets (Tells Make these are commands, not actual filenames)
# ====================================================================================
.PHONY: all build run clean tidy

# Default target when running 'make' without arguments
all: tidy fmt build

## build: Compiles the application binary
build:
	@echo "Building binary to $(BUILD_DIR)/$(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

## run: Builds and runs the local application
run:
	@echo "Running $(BUILD_DIR)/$(BINARY_NAME)..."
	@go run $(MAIN_PACKAGE)

## clean: Removes the compiled binaries and caches
clean:
	@echo "Cleaning up build directory..."
	@rm -rf $(BUILD_DIR)
	@go clean

## fmt: Automatically formats all Go files
fmt:
	@echo "Formatting source files..."
	@go fmt ./...

## tidy: Cleans up and downloads missing Go modules
tidy:
	@echo "Tidying up Go modules..."
	@go mod tidy


