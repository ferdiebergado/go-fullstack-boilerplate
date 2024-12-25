package logging

import (
	"log/slog"
	"os"
)

func handler() slog.Handler {
	logLevel := new(slog.LevelVar)

	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     logLevel,
	}

	var handler slog.Handler

	if os.Getenv("APP_ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		logLevel.Set(slog.LevelDebug)
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	return handler
}

func Init() {
	handler := handler()
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
