package container

import (
	"log/slog"

	"github.com/gerald-lbn/refrain/internal/config"
)

type Container struct {
	Config *config.Config
	Logger *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) *Container {
	return &Container{
		Config: cfg,
		Logger: logger,
	}
}
