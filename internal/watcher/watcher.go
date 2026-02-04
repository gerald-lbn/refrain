package watcher

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

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
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	return filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			w.logger.Debug("Watching directory", "path", path)
			return w.fs.Add(path)
		}
		return nil
	})
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

			if event.Has(fsnotify.Create) {
				st, err := os.Stat(event.Name)
				if err == nil {
					if st.IsDir() && !slices.Contains(w.fs.WatchList(), event.Name) {
						if err := w.Add(event.Name); err != nil {
							w.logger.Error("Failed to add new directory", "path", event.Name, "error", err)
						}
					}
				}
			}
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
