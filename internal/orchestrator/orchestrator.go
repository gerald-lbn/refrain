package orchestrator

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/gerald-lbn/refrain/internal/config"
	"github.com/gerald-lbn/refrain/internal/domain"
	"github.com/gerald-lbn/refrain/internal/helper"
	"github.com/gerald-lbn/refrain/internal/scanner"
)

type Orchestrator struct {
	cfg      *config.Config
	scanner  *scanner.Scanner
	provider domain.LyricsProvider
	logger   *slog.Logger
}

func New(cfg *config.Config, s *scanner.Scanner, p domain.LyricsProvider, logger *slog.Logger) *Orchestrator {
	return &Orchestrator{
		cfg:      cfg,
		scanner:  s,
		provider: p,
		logger:   logger,
	}
}

func (o *Orchestrator) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	for _, lib := range o.cfg.Libraries {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			o.scanLibrary(ctx, path)
		}(lib.Path)
	}

	wg.Wait()
	return nil
}

func (o *Orchestrator) scanLibrary(ctx context.Context, path string) {
	o.logger.InfoContext(ctx, "Starting scan", "path", path)
	tracks, err := o.scanner.Scan(ctx, path)
	if err != nil {
		o.logger.ErrorContext(ctx, "Scan failed", "path", path, "error", err)
		return
	}

	// Worker pool to process tracks
	workers := o.cfg.App.Workers
	var workerWg sync.WaitGroup

	for range workers {
		workerWg.Add(1)
		workerWg.Go(func() {
			defer workerWg.Done()
			for track := range tracks {
				if err := o.processTrack(ctx, track); err != nil {
					o.logger.ErrorContext(ctx, "Failed to process track", "path", track.Path, "error", err)
				}
			}
		})
	}

	workerWg.Wait()
	o.logger.InfoContext(ctx, "Scan complete", "path", path)
}

func (o *Orchestrator) processTrack(ctx context.Context, track domain.Track) error {
	lrcPath := helper.ReplaceExtension(track.Path, ".lrc")
	txtPath := helper.ReplaceExtension(track.Path, ".txt")

	if _, err := os.Stat(lrcPath); err == nil {
		o.logger.DebugContext(ctx, "Lyrics already exist", "path", lrcPath)
		return nil
	}
	if _, err := os.Stat(txtPath); err == nil {
		o.logger.DebugContext(ctx, "Lyrics (txt) already exist", "path", txtPath)
		return nil
	}

	o.logger.DebugContext(ctx, "Searching lyrics", "artist", track.Artist, "album", track.Album, "title", track.Title)

	results, err := o.provider.Search(ctx, track)
	if err != nil {
		return fmt.Errorf("provider search failed: %w", err)
	}

	if len(results) == 0 {
		o.logger.WarnContext(ctx, "No lyrics found", "track", track.Title)
		return nil
	}

	// 3. Save Lyrics (Take the first one, preferably synced)
	bestMatch := results[0]
	for _, l := range results {
		if l.IsSynced {
			bestMatch = l
			break
		}
	}

	savePath := lrcPath
	if !bestMatch.IsSynced {
		savePath = txtPath
	}

	if err := os.WriteFile(savePath, []byte(bestMatch.Text), 0644); err != nil {
		return fmt.Errorf("failed to write lyrics file: %w", err)
	}

	o.logger.DebugContext(ctx, "Saved lyrics", "path", savePath, "source", bestMatch.Source, "synced", bestMatch.IsSynced)
	return nil
}
