package helper

import (
	"path/filepath"
	"strings"
)

var MusicExtensions = map[string]bool{
	".mp3":  true,
	".flac": true,
	".ogg":  true,
	".m4a":  true,
	".wav":  true,
	".aiff": true,
}

// IsMusicFile checks if the file has a supported music extension.
func IsMusicFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return MusicExtensions[ext]
}
