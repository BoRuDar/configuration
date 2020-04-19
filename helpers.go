package configuration

import "os"

type Logger func(format string, v ...interface{})

var (
	gLoggingEnabled  bool
	gFailIfCannotSet bool

	gLogger Logger
)

func logf(format string, args ...interface{}) {
	if gLoggingEnabled {
		gLogger(format, args...)
	}
}

func fatalf(format string, args ...interface{}) {
	if gFailIfCannotSet {
		gLogger(format, args...)
		os.Exit(1)
	}
}
