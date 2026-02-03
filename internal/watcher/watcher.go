package watcher

import (
	"context"
	"log/slog"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	logger *slog.Logger
	fs     *fsnotify.Watcher
}

func New(logger *slog.Logger) (*Watcher, error) {
	fs, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		logger: logger,
		fs:     fs,
	}, nil
}

func (w *Watcher) Add(path string) error {
	w.logger.Info("Watching directory", "path", path)
	return w.fs.Add(path)
}

func (w *Watcher) Run(ctx context.Context) error {
	w.logger.Info("Starting file watcher")
	defer w.fs.Close()

	for {
		select {
		case event, ok := <-w.fs.Events:
			if !ok {
				return nil
			}
			w.logger.Debug("File event", "name", event.Name, "op", event.Op)
			// TODO: Dispatch event to tagger
		case err, ok := <-w.fs.Errors:
			if !ok {
				return nil
			}
			w.logger.Error("Watcher error", "error", err)
		case <-ctx.Done():
			w.logger.Info("Stopping file watcher")
			return ctx.Err()
		}
	}
}
