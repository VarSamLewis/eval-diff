.PHONY: all build test fmt vet lint clean

# Binary name and build settings
BINARY_NAME=bin/eval-diff
GO_FILES=$(shell find . -name '*.go' -type f)

all: fmt vet test build

## build: Compile static Linux/amd64 binary for container execution
build:
	@echo "==> Building static binary..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY_NAME) ./src

## test: Run unit tests
test:
	@echo "==> Running tests..."
	go test -v ./...

## fmt: Format Go source files
fmt:
	@echo "==> Formatting code..."
	gofmt -w .

## vet: Run static analysis
vet:
	@echo "==> Running go vet..."
	go vet ./...

## clean: Remove built binaries
clean:
	@echo "==> Cleaning build artifacts..."
	rm -rf bin/
