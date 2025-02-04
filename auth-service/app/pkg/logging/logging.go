package logging

import (
	"log"
	"log/slog"
	"os"
)

const (
	levelDebug = "debug"
	levelInfo  = "info"
)

func NewLogger(level string) *slog.Logger {
	var logger *slog.Logger

	switch level {
	case levelDebug:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case levelInfo:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	default:
		log.Fatalf("unknown type: %s", level)
	}

	return logger
}
