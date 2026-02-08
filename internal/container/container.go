package container

import (
	"log/slog"

	"github.com/gerald-lbn/refrain/internal/config"
	"github.com/gerald-lbn/refrain/internal/logger"
	"github.com/gerald-lbn/refrain/internal/metadata"
	"github.com/gerald-lbn/refrain/internal/metadata/taglib"
	"github.com/gerald-lbn/refrain/internal/scanner"
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

	return c
}
