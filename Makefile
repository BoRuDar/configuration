COVERAGE_FILE="coverage.out"

.PHONY: test
test:
	go test -v -cover -coverprofile=$(COVERAGE_FILE) -covermode=atomic -count 1  ./...

.PHONY: coverage
coverage: test
	go tool cover -html=$(COVERAGE_FILE)

.PHONY: tools
tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2

linter:
	golangci-lint run ./...
