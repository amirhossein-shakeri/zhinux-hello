SHELL := /bin/sh

GO ?= go
SERVICE_NAME ?= zhinux-hello
BIN_DIR ?= bin
BIN ?= $(BIN_DIR)/$(SERVICE_NAME)
TOOLS_DIR ?= $(CURDIR)/bin/tools
GOLANGCI_LINT ?= $(TOOLS_DIR)/golangci-lint
GOLANGCI_LINT_VERSION ?= v2.4.0
GOLANGCI_LINT_IMAGE ?= golangci/golangci-lint:$(GOLANGCI_LINT_VERSION)
LINT_AUTO_INSTALL ?= true
GO_ENV ?= GOCACHE=$(CURDIR)/.cache/go-build GOMODCACHE=$(CURDIR)/.cache/go-mod

.PHONY: fmt lint test test-integration run build docker-build

fmt:
	$(GO_ENV) $(GO) fmt ./...

lint:
	@set -e; \
	if command -v golangci-lint >/dev/null 2>&1; then \
		$(GO_ENV) golangci-lint run ./...; \
		exit 0; \
	fi; \
	if [ -x "$(GOLANGCI_LINT)" ]; then \
		$(GO_ENV) "$(GOLANGCI_LINT)" run ./...; \
		exit 0; \
	fi; \
	if [ "$(LINT_AUTO_INSTALL)" = "true" ]; then \
		echo "golangci-lint not found; installing $(GOLANGCI_LINT_VERSION) into $(TOOLS_DIR)"; \
		mkdir -p "$(TOOLS_DIR)"; \
		GOBIN="$(TOOLS_DIR)" $(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) || true; \
	fi; \
	if [ -x "$(GOLANGCI_LINT)" ]; then \
		$(GO_ENV) "$(GOLANGCI_LINT)" run ./...; \
		exit 0; \
	fi; \
	if command -v docker >/dev/null 2>&1 && docker image inspect "$(GOLANGCI_LINT_IMAGE)" >/dev/null 2>&1; then \
		docker run --rm -v "$(CURDIR)":/workspace -w /workspace "$(GOLANGCI_LINT_IMAGE)" golangci-lint run ./...; \
		exit 0; \
	fi; \
	echo "golangci-lint is unavailable. Install it in PATH, enable network for auto-install, or pre-pull $(GOLANGCI_LINT_IMAGE)."; \
	exit 2

test:
	$(GO_ENV) $(GO) test ./...

test-integration:
	$(GO_ENV) $(GO) test ./tests/integration/...

run:
	$(GO_ENV) $(GO) run ./cmd

build:
	@mkdir -p $(BIN_DIR)
	$(GO_ENV) $(GO) build -o $(BIN) ./cmd

docker-build:
	docker build -f Dockerfile -t $(SERVICE_NAME):dev ..
