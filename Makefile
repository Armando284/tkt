# tkt - Makefile
.PHONY: build install clean test run help

# Variables
BINARY_NAME=tkt
MAIN_PATH=cmd/tkt/main.go
VERSION=0.1.0-dev
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

# Default target
all: build

# Build the binary
build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Built $(BINARY_NAME) for $(GOOS)/$(GOARCH)"

# Install globally in ~/bin
install:
	go install -ldflags="-s -w -X main.version=$(VERSION)" $(MAIN_PATH)
	@echo "✅ Installed tkt to \$$GOPATH/bin (usually ~/go/bin)"
	@echo "   Make sure ~/go/bin is in your PATH"

# Build a static binary in bin/ (recommended for WSL)
build-static:
	mkdir -p bin
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/$(BINARY_NAME)-$(GOOS)-$(GOARCH) $(MAIN_PATH)
	@echo "✅ Static binary created: bin/$(BINARY_NAME)-$(GOOS)-$(GOARCH)"

# Clean build artifacts
clean:
	rm -rf bin/
	go clean
	@echo "🧹 Cleaned build artifacts"

# Run directly without building
run:
	go run $(MAIN_PATH)

# Run with hello command (quick test)
hello:
	go run $(MAIN_PATH) hello Armando

# Run tests (when we have them)
test:
	go test ./...

# Show help
help:
	@echo "Available commands:"
	@echo "  make build          - Build binary into bin/"
	@echo "  make install        - Install to GOPATH/bin"
	@echo "  make build-static   - Build static binary"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make run            - Run directly with go run"
	@echo "  make hello          - Quick test: tkt hello Armando"
	@echo "  make help           - Show this help"
