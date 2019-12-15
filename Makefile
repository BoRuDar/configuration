APP_NAME="PlayersProfile"
COVERAGE_FILE="coverage.txt"

test:
	go test -v -cover -coverprofile=$(COVERAGE_FILE)  ./...

coverage: test
	go tool cover -html=$(COVERAGE_FILE)