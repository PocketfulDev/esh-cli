# Makefile for esh-cli

BINARY_NAME=esh-cli
MAIN_PACKAGE=.
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: build clean install test test-coverage test-coverage-json test-coverage-check test-race test-all test-ci fmt vet check help release-build

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Build binaries for all platforms
release-build:
	@echo "Building release binaries..."
	@mkdir -p dist
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	go clean
	rm -f $(BINARY_NAME)
	rm -rf dist/

# Install the binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	go install $(MAIN_PACKAGE)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -v -race ./...

# Run tests with coverage and JSON output (matches GitHub Actions)
test-coverage-json:
	@echo "Running tests with coverage and JSON output..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic -json ./... > test-results.json
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out > coverage-func.txt
	@echo "Coverage files generated:"
	@echo "  - coverage.out (Go coverage profile)"
	@echo "  - coverage.html (HTML report)"
	@echo "  - coverage-func.txt (Function breakdown)"
	@echo "  - test-results.json (JSON test output)"

# Check coverage thresholds (matches GitHub Actions)
test-coverage-check: test-coverage
	@echo "Checking coverage thresholds..."
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Total coverage: $${COVERAGE}%"; \
	if command -v bc >/dev/null 2>&1; then \
		if [ $$(echo "$${COVERAGE} < 30" | bc -l) -eq 1 ]; then \
			echo "⚠️  Warning: Coverage below 30% threshold"; \
		else \
			echo "✅ Coverage meets 30% threshold"; \
		fi; \
	else \
		echo "ℹ️  bc not available, skipping threshold check"; \
	fi; \
	UTILS_COVERAGE=$$(go test -coverprofile=utils_coverage.out -covermode=atomic ./pkg/utils 2>/dev/null && go tool cover -func=utils_coverage.out | grep total | awk '{print $$3}' | sed 's/%//' || echo "0"); \
	echo "Utils package coverage: $${UTILS_COVERAGE}%"; \
	if command -v bc >/dev/null 2>&1 && [ -n "$${UTILS_COVERAGE}" ] && [ "$${UTILS_COVERAGE}" != "0" ]; then \
		if [ $$(echo "$${UTILS_COVERAGE} < 60" | bc -l) -eq 1 ]; then \
			echo "⚠️  Warning: Utils coverage below 60% threshold"; \
		else \
			echo "✅ Utils coverage meets 60% threshold"; \
		fi; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run CI-style testing (matches GitHub Actions exactly)
test-ci: test-coverage-json test-coverage-check
	@echo "✅ CI-style testing complete"

# Run comprehensive tests (coverage + race detection)
test-all: test-coverage test-race

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Run all checks
check: fmt vet test-all

# Show help
help:
	@echo "Available targets:"
	@echo "  build               - Build the binary"
	@echo "  release-build       - Build binaries for all platforms"
	@echo "  clean               - Clean build artifacts"
	@echo "  install             - Install the binary"
	@echo "  test                - Run basic tests"
	@echo "  test-coverage       - Run tests with coverage report"
	@echo "  test-coverage-json  - Run tests with coverage and JSON output (CI-style)"
	@echo "  test-coverage-check - Check coverage meets thresholds"
	@echo "  test-race           - Run tests with race detection"
	@echo "  test-all            - Run comprehensive tests (coverage + race)"
	@echo "  test-ci             - Run full CI-style testing"
	@echo "  fmt                 - Format code"
	@echo "  vet                 - Vet code"
	@echo "  check               - Run fmt, vet, and comprehensive tests"
	@echo "  help                - Show this help"
