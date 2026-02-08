package domain

import "context"

type LyricsProvider interface {
	Name() string
	Search(ctx context.Context, track Track) ([]Lyrics, error)
	Get(ctx context.Context, id string) (Lyrics, error)
}

type LibraryScanner interface {
	Scan(ctx context.Context, path string) (<-chan Track, error)
}
