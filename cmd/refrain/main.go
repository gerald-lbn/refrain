package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gerald-lbn/refrain/internal/config"
	"github.com/gerald-lbn/refrain/internal/container"
	"github.com/gerald-lbn/refrain/internal/logger"
	"github.com/gerald-lbn/refrain/internal/watcher"
)

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, stop := mainContext(context.Background())
	defer stop()

	cfg, err := config.Load("/config/config.yaml")
	if err != nil {
		l.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	l = logger.New(cfg)

	// Initialize Watcher
	w, err := watcher.New(l)
	if err != nil {
		l.Error("Failed to initialize watcher", "error", err)
		os.Exit(1)
	}

	// Add configured libraries to watcher
	for _, lib := range cfg.Libraries {
		if err := w.Add(lib.Path); err != nil {
			l.Error("Failed to watch library", "name", lib.Name, "path", lib.Path, "error", err)
		}
	}

	c := container.New(cfg, l, w)

	l.Info("Refrain is running...")
	if err := c.Run(ctx); err != nil && err != context.Canceled {
		l.Error("Container stopped with error", "error", err)
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
