package domain

import (
	"context"
	"time"
)

type LyricsProvider interface {
	Search(ctx context.Context, track Track) ([]Lyrics, error)
	Download(ctx context.Context, id string) (Lyrics, error)
}

type LibraryScanner interface {
	Scan(ctx context.Context, path string) (<-chan Track, error)
}

type Scheduler interface {
	AddFunc(interval time.Duration, cmd func())
	Start(ctx context.Context)
	Stop()
}
