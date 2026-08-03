CHMOD_CMD = chmod +x .githooks/pre-commit
ifeq ($(OS),Windows_NT)
    CHMOD_CMD = echo "Skipping chmod on Windows"
endif

.PHONY: install-hooks
install-hooks:
	@echo "Installing git hooks from .githooks directory..."
	@$(CHMOD_CMD)
	@git config core.hooksPath .githooks

.PHONY: help
.DEFAULT_GOAL := help
help:
	@echo "Available commands:"
	@echo "  make format          Format code"
	@echo "  make lint            Check code style"
	@echo "  make typecheck       Run type checks"
	@echo "  make test            Run tests with coverage"
	@echo "  make test-fast       Run tests without coverage"
	@echo "  make coverage-html   Generate HTML coverage reports"
	@echo "  make checks          Run lint, typecheck, and test"

.PHONY: format
format:
	$(MAKE) -C python format
	$(MAKE) -C go format

.PHONY: lint
lint:
	$(MAKE) -C python lint
	$(MAKE) -C go lint

.PHONY: typecheck
typecheck:
	$(MAKE) -C python typecheck
	$(MAKE) -C go typecheck

.PHONY: build
build:
	$(MAKE) -C go build

.PHONY: test
test:
	$(MAKE) -C python test
	$(MAKE) -C go test

.PHONY: test-fast
test-fast:
	$(MAKE) -C python test-fast
	$(MAKE) -C go test-fast

.PHONY: coverage-html
coverage-html:
	$(MAKE) -C python coverage-html
	$(MAKE) -C go coverage-html

.PHONY: checks
checks:
	$(MAKE) -C python checks
	$(MAKE) -C go checks
