package helper

import (
	"path/filepath"
	"strings"
)

// ReplaceExtension replaces the extension of the given path with newExt.
// newExt should include the dot.
func ReplaceExtension(path, newExt string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext) + newExt
}
