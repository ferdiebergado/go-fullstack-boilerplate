package logging

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func New() *Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{AddSource: false}

	if os.Getenv("ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &Logger{slog.New(handler)}
}

func (l *Logger) Fatal(reason string, err error) {
	l.Logger.Error(
		"Fatal error occurred",
		slog.String("reason", reason),
		slog.String("error", err.Error()),
		slog.String("severity", "FATAL"),
	)

	os.Exit(1)
}
