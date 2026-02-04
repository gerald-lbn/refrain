package container

import (
	"context"
	"log/slog"

	"github.com/gerald-lbn/refrain/internal/config"
	"github.com/gerald-lbn/refrain/internal/watcher"
	"golang.org/x/sync/errgroup"
)

type Container struct {
	Config  *config.Config
	Logger  *slog.Logger
	Watcher *watcher.Watcher
}

func New(cfg *config.Config, logger *slog.Logger, w *watcher.Watcher) *Container {
	return &Container{
		Config:  cfg,
		Logger:  logger,
		Watcher: w,
	}
}

func (c *Container) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// Start Watcher
	if c.Watcher != nil {
		g.Go(func() error {
			return c.Watcher.Run(ctx)
		})
	}

	return g.Wait()
}
