package domain

import "time"

type Track struct {
	Path     string
	Title    string
	Artist   string
	Album    string
	Duration time.Duration
}

type Lyrics struct {
	Text         string
	IsSynced     bool
	Instrumental bool
	Source       string
}
