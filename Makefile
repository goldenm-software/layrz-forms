.PHONY: help
.DEFAULT_GOAL := help

help:
	@echo "Available targets:"
	@echo "  make format          Format code"
	@echo "  make lint            Check code style"
	@echo "  make typecheck       Run type checks"
	@echo "  make test            Run all tests with coverage (90% threshold enforced)"
	@echo "  make test-fast       Run tests without coverage instrumentation"
	@echo "  make coverage-html   Generate HTML coverage reports"
	@echo "  make checks          Run lint, typecheck, and tests (pre-push gate)"

.PHONY: format
format:
	$(MAKE) -C python format

.PHONY: lint
lint:
	$(MAKE) -C python lint

.PHONY: typecheck
typecheck:
	$(MAKE) -C python typecheck

.PHONY: test
test:
	$(MAKE) -C python test

.PHONY: test-fast
test-fast:
	$(MAKE) -C python test-fast

.PHONY: coverage-html
coverage-html:
	$(MAKE) -C python coverage-html

.PHONY: checks
checks:
	$(MAKE) -C python checks