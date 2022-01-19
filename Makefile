COVERAGE_FILE="coverage.out"

test:
	go test -v -cover -coverprofile=$(COVERAGE_FILE) -covermode=atomic -count 1  ./...

coverage: test
	go tool cover -html=$(COVERAGE_FILE)
