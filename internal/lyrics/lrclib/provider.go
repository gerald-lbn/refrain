package lrclib

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/gerald-lbn/refrain/internal/lyrics"
)

const (
	API_URL = "https://lrclib.net/api"
)

type Provider struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Provider {
	return &Provider{
		logger: logger,
	}
}

type LrcLibResult struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	TrackName    string  `json:"trackName"`
	ArtistName   string  `json:"artistName"`
	AlbumName    string  `json:"albumName"`
	Duration     float64 `json:"duration"`
	Instrumental bool    `json:"instrumental"`
	PlainLyrics  string  `json:"plainLyrics"`
	SyncedLyrics string  `json:"syncedLyrics"`
}

func lrcLibResultToLyricsResult(r LrcLibResult) lyrics.Result {
	return lyrics.Result{
		Id:           r.Id,
		Artist:       r.ArtistName,
		Title:        r.TrackName,
		Album:        r.AlbumName,
		Instrumental: r.Instrumental,
		Duration:     r.Duration,
	}
}

func lrcLibResultToLyrics(r LrcLibResult) lyrics.Lyrics {
	return lyrics.Lyrics{
		Synced: r.SyncedLyrics,
		Plain:  r.PlainLyrics,
	}
}

func (p *Provider) Search(ctx context.Context, search lyrics.Search) ([]lyrics.Result, error) {
	params := url.Values{}
	if search.Query != nil {
		params.Add("query", *search.Query)
	} else {
		if search.Track != nil {
			params.Add("track_name", *search.Track)
		}
		if search.Artist != nil {
			params.Add("artist_name", *search.Artist)
		}
		if search.Album != nil {
			params.Add("album_name", *search.Album)
		}
		if search.Duration != nil {
			params.Add("duration", fmt.Sprintf("%f", *search.Duration))
		}
	}
	url := fmt.Sprintf("%s/search?%s", API_URL, params.Encode())
	p.logger.Debug("fetching", slog.String("url", url))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var lrcLibResults []LrcLibResult
	err = json.Unmarshal(body, &lrcLibResults)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results []lyrics.Result
	for _, result := range lrcLibResults {
		results = append(results, lrcLibResultToLyricsResult(result))
	}

	return results, nil
}

func (p *Provider) GetLyrics(ctx context.Context, id int) (*lyrics.Lyrics, error) {
	url := fmt.Sprintf("%s/get/%d", API_URL, id)
	p.logger.Debug("fetching", slog.String("url", url))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result LrcLibResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	lyrics := lrcLibResultToLyrics(result)
	return &lyrics, nil
}
