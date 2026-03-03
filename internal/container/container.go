package container

import (
	"log/slog"

	"github.com/gerald-lbn/refrain/internal/config"
	"github.com/gerald-lbn/refrain/internal/domain"
	"github.com/gerald-lbn/refrain/internal/logger"
	"github.com/gerald-lbn/refrain/internal/metadata"
	"github.com/gerald-lbn/refrain/internal/metadata/taglib"
	"github.com/gerald-lbn/refrain/internal/orchestrator"
	"github.com/gerald-lbn/refrain/internal/provider/lrclib"
	"github.com/gerald-lbn/refrain/internal/scanner"
	"github.com/gerald-lbn/refrain/internal/watcher"
	"go.uber.org/dig"
)

func Build() *dig.Container {
	c := dig.New()

	c.Provide(func() *config.Config {
		return config.Load()
	})

	c.Provide(func(cfg *config.Config) *slog.Logger {
		return logger.New(cfg.LogLevel)
	})

	c.Provide(func() metadata.Reader {
		return taglib.New()
	})

	c.Provide(scanner.New)

	c.Provide(func(logger *slog.Logger) domain.LyricsProvider {
		return lrclib.New(logger, nil)
	})

	c.Provide(func(logger *slog.Logger) (*watcher.Watcher, error) {
		return watcher.New(logger.With("component", "watcher"))
	})

	c.Provide(orchestrator.New)

	return c
}
