package metadata

import "github.com/gerald-lbn/refrain/internal/domain"

//go:generate mockgen -destination=../mocks/metadata/mock_reader.go -package=mock_metadata github.com/gerald-lbn/refrain/internal/metadata Reader

type Reader interface {
	Read(path string) (*domain.Track, error)
}
