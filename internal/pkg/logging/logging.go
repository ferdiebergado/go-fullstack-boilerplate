package logging

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func New() *Logger {
	handler := getHandler()
	return &Logger{slog.New(handler)}
}

func getHandler() slog.Handler {
	var handler slog.Handler

	opts := &slog.HandlerOptions{AddSource: false}

	if os.Getenv("ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	return handler
}

func SetLogger() {
	handler := getHandler()
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func (l *Logger) Fatal(reason string, err error) {
	l.Logger.Error(
		"Fatal error occurred",
		"reason", reason,
		"error", err.Error(),
		"severity", "FATAL",
	)

	panic(reason)
}
