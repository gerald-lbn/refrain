package watcher

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/gerald-lbn/refrain/internal/helper"
)

// Watcher watches directories for new music files using fsnotify.
type Watcher struct {
	logger  *slog.Logger
	watcher *fsnotify.Watcher
}

// New creates a new Watcher.
func New(logger *slog.Logger) (*Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		logger:  logger,
		watcher: fw,
	}, nil
}

// Watch starts watching the given root paths (recursively) for new music files.
// It returns a channel of file paths for newly created music files.
// The channel is closed when the context is cancelled or an unrecoverable error occurs.
func (w *Watcher) Watch(ctx context.Context, paths []string) (<-chan string, error) {
	for _, root := range paths {
		if err := w.addRecursive(root); err != nil {
			return nil, err
		}
	}

	out := make(chan string)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}
				w.handleEvent(ctx, event, out)
			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				w.logger.ErrorContext(ctx, "Watcher error", "error", err)
			}
		}
	}()

	return out, nil
}

// Close stops the underlying fsnotify watcher.
func (w *Watcher) Close() error {
	return w.watcher.Close()
}

// handleEvent processes a single fsnotify event.
func (w *Watcher) handleEvent(ctx context.Context, event fsnotify.Event, out chan<- string) {
	if !event.Has(fsnotify.Create) && !event.Has(fsnotify.Rename) {
		return
	}

	// If a new directory was created, watch it recursively.
	info, err := os.Stat(event.Name)
	if err != nil {
		return
	}

	if info.IsDir() {
		w.logger.DebugContext(ctx, "New directory detected, adding watch", "path", event.Name)
		if err := w.addRecursive(event.Name); err != nil {
			w.logger.ErrorContext(ctx, "Failed to watch new directory", "path", event.Name, "error", err)
		}
		return
	}

	if !helper.IsMusicFile(event.Name) {
		return
	}

	w.logger.DebugContext(ctx, "New music file detected", "path", event.Name)

	select {
	case out <- event.Name:
	case <-ctx.Done():
	}
}

// addRecursive adds a watch on the given directory and all its subdirectories.
func (w *Watcher) addRecursive(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			w.logger.Warn("Error walking path during watch setup", "path", path, "error", err)
			return nil
		}

		if !d.IsDir() {
			return nil
		}

		w.logger.Debug("Watching directory", "path", path)
		return w.watcher.Add(path)
	})
}
