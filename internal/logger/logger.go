package logger

import (
	"log/slog"
	"os"

	"github.com/gerald-lbn/refrain/internal/config"
)

func New(cfg *config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: stringToLevel(cfg.Logger.Level),
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

func stringToLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
