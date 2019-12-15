APP_NAME="PlayersProfile"
COVERAGE_FILE="coverage.out"

test:
	go test -v -cover -coverprofile=$(COVERAGE_FILE) -covermode=atomic  ./...

coverage: test
	go tool cover -html=$(COVERAGE_FILE)