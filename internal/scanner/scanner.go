package scanner

import (
	"context"
	"io/fs"
	"log/slog"
	"path/filepath"

	"github.com/gerald-lbn/refrain/internal/domain"
	"github.com/gerald-lbn/refrain/internal/helper"
	"github.com/gerald-lbn/refrain/internal/metadata"
)

// Scanner scans a directory for music files.
type Scanner struct {
	logger         *slog.Logger
	metadataReader metadata.Reader
}

// New creates a new Scanner.
func New(logger *slog.Logger, metadataReader metadata.Reader) *Scanner {
	return &Scanner{
		logger:         logger,
		metadataReader: metadataReader,
	}
}

// Scan walks the given path recursively and sends found tracks to the returned channel.
func (s *Scanner) Scan(ctx context.Context, path string) (<-chan domain.Track, error) {
	out := make(chan domain.Track)

	go func() {
		defer close(out)

		err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				s.logger.ErrorContext(ctx, "Error walking path", "path", path, "error", err)
				return nil
			}

			if d.IsDir() {
				return nil
			}

			if !helper.IsMusicFile(path) {
				return nil
			}

			// Check context cancellation
			select {
			case <-ctx.Done():
				return filepath.SkipAll
			default:
			}

			tags, err := s.metadataReader.Read(path)
			if err != nil {
				s.logger.WarnContext(ctx, "Failed to read tags", "path", path, "error", err)
				return nil
			}

			track := domain.Track{
				Path:     path,
				Title:    tags.Title,
				Artist:   tags.Artist,
				Album:    tags.Album,
				Duration: tags.Duration,
			}

			// Send to channel
			select {
			case out <- track:
			case <-ctx.Done():
				return filepath.SkipAll
			}

			return nil
		})

		if err != nil {
			s.logger.ErrorContext(ctx, "WalkDir failed", "error", err)
		}
	}()

	return out, nil
}
