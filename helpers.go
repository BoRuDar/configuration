package configuration

import (
	"log"
	"os"
)

var (
	gLoggingEnabled  bool
	gFailIfCannotSet bool
)

func logf(format string, args ...interface{}) {
	if gLoggingEnabled {
		log.Printf(format, args...)
	}
}

func fail() {
	if gFailIfCannotSet {
		os.Exit(1)
	}
}
