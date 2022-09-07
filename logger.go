package hooks

// Logger is a function to handle logging
type Logger func(format string, args ...any)

// logger stores the logging function
var logger Logger

// SetLogger sets a logger function to log hook events
func SetLogger(function Logger) {
	logger = function
}

// logf logs output to the logger
func logf(format string, args ...any) {
	if logger == nil {
		return
	}

	logger(format, args...)
}
