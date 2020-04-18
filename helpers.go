package configuration

import "os"

type Logger func(format string, v ...interface{})

var (
	gLoggingEnabled  bool
	gFailIfCannotSet bool

	logger Logger
)

func logf(format string, args ...interface{}) {
	if gLoggingEnabled {
		logger(format, args...)
	}
}

func fatalf(format string, args ...interface{}) {
	if gFailIfCannotSet {
		logger(format, args...)
		os.Exit(1)
	}
}
