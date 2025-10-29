.PHONY: all build clean install test

# Binary output directory
BIN_DIR := bin
CLI_BINARY := $(BIN_DIR)/bsubio

# Go command
GO := go
GOFLAGS := -v

# Build the CLI
all: build

build: $(CLI_BINARY)

$(CLI_BINARY):
	@echo "Building bsubio CLI..."
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(CLI_BINARY) ./cmd/bsubio

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)

# Install the CLI to system PATH
install: build
	@echo "Installing bsubio to /usr/local/bin..."
	@install -m 0755 $(CLI_BINARY) /usr/local/bin/bsubio
	@echo "Installation complete!"

# Run tests
test:
	$(GO) test -v ./...

# Download dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy

# Show help
help:
	@echo "Makefile targets:"
	@echo "  build    - Build the CLI binary"
	@echo "  clean    - Remove build artifacts"
	@echo "  install  - Install CLI to /usr/local/bin"
	@echo "  test     - Run tests"
	@echo "  deps     - Download and tidy dependencies"
	@echo "  help     - Show this help message"
