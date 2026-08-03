.PHONY: help
.DEFAULT_GOAL := help

help:
	@echo "Available targets:"
	@echo "  make format          Format code with gofmt"
	@echo "  make lint            Check code style with gofmt"
	@echo "  make typecheck       Run type checks with go vet"
	@echo "  make build           Build Go binaries"
	@echo "  make test            Run tests with coverage (enforces 90% threshold)"
	@echo "  make test-fast       Run tests without coverage instrumentation"
	@echo "  make coverage-html   Generate HTML coverage report in coverage.html"
	@echo "  make checks          Run lint, typecheck, build, and test (pre-push gate)"

.PHONY: format
format:
	gofmt -w .

.PHONY: lint
lint:
	@test -z "$$(gofmt -l .)" || { gofmt -l .; exit 1; }

.PHONY: typecheck
typecheck:
	go vet ./...

.PHONY: build
build:
	go build ./...

# Current coverage: 90.6% | Matches the Python threshold of 90%
COVERAGE_THRESHOLD := 90

.PHONY: test
test:
	go test -coverprofile=coverage.out ./...
	@total=$$(go tool cover -func=coverage.out | awk '/^total:/ {print $$3}' | tr -d '%'); \
	echo "total coverage: $$total%"; \
	awk -v t="$$total" -v min=$(COVERAGE_THRESHOLD) 'BEGIN { if (t+0 < min+0) { printf "FAIL: coverage %s%% is below %s%%\n", t, min; exit 1 } }'

.PHONY: test-fast
test-fast:
	go test ./...

.PHONY: coverage-html
coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "HTML coverage report generated in coverage.html"

.PHONY: checks
checks: lint typecheck build test
