package container

import (
	"log/slog"

	"github.com/gerald-lbn/refrain/internal/config"
	scheduler "github.com/gerald-lbn/refrain/internal/cron"
	"github.com/gerald-lbn/refrain/internal/domain"
	"github.com/gerald-lbn/refrain/internal/logger"
	"github.com/gerald-lbn/refrain/internal/metadata"
	"github.com/gerald-lbn/refrain/internal/metadata/taglib"
	"github.com/gerald-lbn/refrain/internal/orchestrator"
	"github.com/gerald-lbn/refrain/internal/provider/lrclib"
	"github.com/gerald-lbn/refrain/internal/scanner"
	goCron "github.com/robfig/cron/v3"
	"go.uber.org/dig"
)

func Build() *dig.Container {
	c := dig.New()

	c.Provide(func() (*config.Config, error) {
		return config.LoadConfig(config.DefaultConfigPath)
	})

	c.Provide(func(cfg *config.Config) *slog.Logger {
		return logger.New(cfg.Log.Level)
	})

	c.Provide(func() metadata.Reader {
		return taglib.New()
	})

	c.Provide(scanner.New)

	c.Provide(func(logger *slog.Logger) domain.LyricsProvider {
		return lrclib.New(logger, nil)
	})

	c.Provide(func(logger *slog.Logger) domain.Scheduler {
		return scheduler.New(logger.With("component", "scheduler"), goCron.New())
	})

	c.Provide(orchestrator.New)

	return c
}
