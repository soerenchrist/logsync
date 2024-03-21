package log

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func New(level slog.Level) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger = slog.New(handler)
	logger.Info("Setting log level", "level", level)
	slog.SetDefault(logger)
	return logger
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}
