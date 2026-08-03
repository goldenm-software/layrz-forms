.PHONY: help
.DEFAULT_GOAL := help

help:
	@echo "Available targets:"
	@echo "  make format          Format code with ruff"
	@echo "  make lint            Check code style with ruff"
	@echo "  make typecheck       Run type checks with ty"
	@echo "  make test            Run tests with coverage (enforces 80% threshold)"
	@echo "  make test-fast       Run tests without coverage instrumentation"
	@echo "  make coverage-html   Generate HTML coverage report in htmlcov/"
	@echo "  make checks          Run lint, typecheck, and test (pre-push gate)"

.PHONY: format
format:
	uv run ruff format .

.PHONY: lint
lint:
	uv run ruff check

.PHONY: typecheck
typecheck:
	uv run ty check

.PHONY: test
test:
	uv run pytest -q --cov=layrz_forms --cov-report=term-missing --cov-fail-under=90

.PHONY: test-fast
test-fast:
	uv run pytest -q

.PHONY: coverage-html
coverage-html:
	uv run pytest --cov=layrz_forms --cov-report=html
	@echo "HTML coverage report generated in htmlcov/"

.PHONY: checks
checks: lint typecheck test