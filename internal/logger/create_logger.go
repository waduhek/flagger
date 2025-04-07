package logger

import (
	"log/slog"
	"os"
)

// CreateLogger creates a new logger.
func CreateLogger() *slog.Logger {
	handler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)
	logger := slog.New(handler)

	return logger
}
