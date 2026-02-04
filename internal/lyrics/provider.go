package lyrics

import "context"

type Lyrics struct {
	Synced string
	Plain  string
}

type Result struct {
	Id           int
	Artist       string
	Title        string
	Album        string
	Instrumental bool
	Duration     float64
}

type Search struct {
	Track    *string
	Artist   *string
	Album    *string
	Query    *string
	Duration *float64
}

type Provider interface {
	Search(ctx context.Context, search Search) ([]Result, error)
	GetLyrics(ctx context.Context, id int) (*Lyrics, error)
}
