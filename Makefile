SHELL := /bin/sh

GO ?= go
SERVICE_NAME ?= zhinux-hello
BIN_DIR ?= bin
BIN ?= $(BIN_DIR)/$(SERVICE_NAME)

.PHONY: fmt lint test test-integration run build docker-build

fmt:
	$(GO) fmt ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint is required for lint"; exit 1; }
	golangci-lint run ./...

test:
	$(GO) test ./...

test-integration:
	$(GO) test ./tests/integration/...

run:
	$(GO) run ./cmd

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN) ./cmd

docker-build:
	docker build -t $(SERVICE_NAME):dev .
