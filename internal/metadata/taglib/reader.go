package taglib

import (
	"github.com/gerald-lbn/refrain/internal/domain"
	t "go.senan.xyz/taglib"
)

type Reader struct{}

func New() *Reader {
	return &Reader{}
}

func (r *Reader) Read(path string) (*domain.Track, error) {
	tags, err := t.ReadTags(path)
	if err != nil {
		return nil, err
	}

	props, err := t.ReadProperties(path)
	if err != nil {
		return nil, err
	}

	track := &domain.Track{
		Duration: props.Length,
	}

	if len(tags[t.Title]) > 0 {
		track.Title = tags[t.Title][0]
	}
	if len(tags[t.AlbumArtist]) > 0 {
		track.Artist = tags[t.AlbumArtist][0]
	}
	if len(tags[t.Album]) > 0 {
		track.Album = tags[t.Album][0]
	}

	return track, nil
}
