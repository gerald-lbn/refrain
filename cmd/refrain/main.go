package main

import (
	"log/slog"
	"os"

	"github.com/gerald-lbn/refrain/internal/config"
	"github.com/gerald-lbn/refrain/internal/container"
	"github.com/gerald-lbn/refrain/internal/logger"
)

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load("/config/config.yaml")
	if err != nil {
		l.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	l = logger.New(cfg)

	c := container.New(cfg, l)

	l.Info("Refrain Backend Initialized")
	if c.Config != nil {
		l.Info("Loaded libraries", "count", len(c.Config.Libraries))
		for _, lib := range c.Config.Libraries {
			l.Info("Library loaded", "name", lib.Name, "path", lib.Path)
		}
	}
}
