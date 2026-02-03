package container

import (
	"github.com/gerald-lbn/refrain/internal/config"
)

type Container struct {
	Config *config.Config
}

func New(cfg *config.Config) *Container {
	return &Container{
		Config: cfg,
	}
}
