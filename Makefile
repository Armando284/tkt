# tkt - Makefile
.PHONY: build build-debug clean run dev

BINARY_NAME=tkt-bin
MAIN_PATH=./cmd/tkt
VERSION=0.1.0-dev

all: build

# Build correcto - compila todo el paquete cmd/tkt
build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Built $(BINARY_NAME) → bin/$(BINARY_NAME)"

build-debug:
	go build -tags=debug -ldflags="-s -w -X main.version=$(VERSION)" -o bin/$(BINARY_NAME)-debug $(MAIN_PATH)
	@echo "✅ Built with DEBUG logs → bin/$(BINARY_NAME)-debug"

clean:
	rm -rf bin/
	go clean -cache
	@echo "🧹 Cleaned"

run:
	go run $(MAIN_PATH)

# Para probar comandos rápido
dev:
	go run $(MAIN_PATH) start 23

# Para probar list
dev-list:
	go run $(MAIN_PATH) list