package logging

import (
	"log/slog"
)

type Logger struct {
	*slog.Logger
}

func New(logger *slog.Logger) *Logger {
	return &Logger{logger}
}
