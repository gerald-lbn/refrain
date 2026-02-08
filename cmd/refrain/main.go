package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gerald-lbn/refrain/internal/config"
	"github.com/gerald-lbn/refrain/internal/container"
	"github.com/gerald-lbn/refrain/internal/scanner"
)

func main() {
	ctx, cancel := mainContext(context.Background())
	defer cancel()

	c := container.Build()

	err := c.Invoke(func(s *scanner.Scanner, cfg *config.Config, logger *slog.Logger) {
		logger.Info("Starting Refrain...")

		for _, lib := range cfg.Libraries {
			logger.Info("Scanning library", "path", lib.Path)
		}

		<-ctx.Done()
		logger.Info("Shutting down...")
	})

	if err != nil {
		slog.Error("Failed to start application", "error", err)
		os.Exit(1)
	}
}

// mainContext returns a context that is cancelled when the process receives a signal to exit.
func mainContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGABRT,
	)
}
