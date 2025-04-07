package logger

type Logger interface {
	// Debug logs a message the debug level.
	Debug(msg string, args ...any)

	// Info logs a message at the information level.
	Info(msg string, args ...any)

	// Warn logs a message at the warning level.
	Warn(msg string, args ...any)

	// Error logs a message the error level.
	Error(msg string, args ...any)
}
