package domain

import "context"

type LyricsProvider interface {
	Search(ctx context.Context, track Track) ([]Lyrics, error)
	Download(ctx context.Context, id string) (Lyrics, error)
}

type LibraryScanner interface {
	Scan(ctx context.Context, path string) (<-chan Track, error)
}
