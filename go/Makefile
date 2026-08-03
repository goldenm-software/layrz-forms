.PHONY: help
.DEFAULT_GOAL := help

help:
	@echo "Available targets:"
	@echo "  make format          Format code with gofmt"
	@echo "  make lint            Check code style with gofmt and golangci-lint"
	@echo "  make typecheck       Run type checks with go vet"
	@echo "  make build           Build Go binaries"
	@echo "  make test            Run tests with coverage (enforces 90% threshold)"
	@echo "  make test-fast       Run tests without coverage instrumentation"
	@echo "  make coverage-html   Generate HTML coverage report in coverage.html"
	@echo "  make checks          Run pre-commit checks (lint, vet, golangci-lint, go mod tidy)"

.PHONY: format
format:
	gofmt -w .

.PHONY: lint
lint:
	@test -z "$$(gofmt -l .)" || { gofmt -l .; exit 1; }
	golangci-lint config verify
	golangci-lint run ./...

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
checks:
	@echo "Running pre-commit checks..."
	go fmt ./...
	go vet ./...
	golangci-lint config verify
	golangci-lint run ./...
	@echo "Checking if go.mod is tidy..."
	@cp go.mod go.mod.tmp
	@if [ -f go.sum ]; then cp go.sum go.sum.tmp; fi
	@go mod tidy
	@if ! cmp -s go.mod go.mod.tmp; then \
		echo "Error: go.mod was modified by 'go mod tidy'."; \
		mv go.mod.tmp go.mod; \
		if [ -f go.sum.tmp ]; then mv go.sum.tmp go.sum; fi; \
		exit 1; \
	fi
	@if [ -f go.sum ]; then \
		if ! cmp -s go.sum go.sum.tmp; then \
			echo "Error: go.sum was modified by 'go mod tidy'."; \
			mv go.mod.tmp go.mod; \
			mv go.sum.tmp go.sum; \
			exit 1; \
		fi; \
		rm go.sum.tmp; \
	fi
	@rm go.mod.tmp
	@echo "All pre-commit checks passed."
