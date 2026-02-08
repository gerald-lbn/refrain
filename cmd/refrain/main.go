package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gerald-lbn/refrain/internal/config"
	"github.com/gerald-lbn/refrain/internal/container"
	"github.com/gerald-lbn/refrain/internal/orchestrator"
)

func main() {
	ctx, cancel := mainContext(context.Background())
	defer cancel()

	c := container.Build()

	err := c.Invoke(func(orc *orchestrator.Orchestrator, cfg *config.Config, logger *slog.Logger) {
		logger.Info("Starting Refrain...")

		if err := orc.Run(context.Background()); err != nil {
			logger.Error("Orchestrator failed", "error", err)
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
