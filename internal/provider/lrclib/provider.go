package lrclib

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gerald-lbn/refrain/internal/domain"
)

const (
	baseURL = "https://lrclib.net/api"
)

type Provider struct {
	client *http.Client
	logger *slog.Logger
}

func New(logger *slog.Logger, client *http.Client) *Provider {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &Provider{
		client: client,
		logger: logger,
	}
}

func (p *Provider) Name() string {
	return "LRCLIB"
}

// Search searches for lyrics for the given track.
func (p *Provider) Search(ctx context.Context, track domain.Track) ([]domain.Lyrics, error) {
	u, err := url.Parse(baseURL + "/search")
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	q := u.Query()
	q.Set("track_name", track.Title)
	q.Set("artist_name", track.Artist)
	q.Set("album_name", track.Album)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var results []lrclibTrack
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var lyrics []domain.Lyrics
	for _, res := range results {
		text := res.SyncedLyrics
		isSynced := true
		if text == "" {
			text = res.PlainLyrics
			isSynced = false
		}

		if text == "" && res.Instrumental {
			text = "Instrumental"
			isSynced = false
		}

		if text == "" && !res.Instrumental {
			continue
		}

		lyrics = append(lyrics, domain.Lyrics{
			Text:         text,
			IsSynced:     isSynced,
			Instrumental: res.Instrumental,
			Source:       "LRCLIB",
		})
	}

	return lyrics, nil
}

// Download gets lyrics by ID.
func (p *Provider) Download(ctx context.Context, id string) (domain.Lyrics, error) {
	u, err := url.Parse(baseURL + "/get/" + id)
	if err != nil {
		return domain.Lyrics{}, fmt.Errorf("failed to parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return domain.Lyrics{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return domain.Lyrics{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.Lyrics{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res lrclibTrack
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return domain.Lyrics{}, fmt.Errorf("failed to decode response: %w", err)
	}

	text := ""
	if res.SyncedLyrics == "" && res.PlainLyrics == "" && !res.Instrumental {
		return domain.Lyrics{}, fmt.Errorf("no lyrics found for id: %s", id)
	} else if res.Instrumental {
		text = "Instrumental"
	} else if res.SyncedLyrics != "" {
		text = res.SyncedLyrics
	} else if res.PlainLyrics != "" {
		text = res.PlainLyrics
	}
	isSynced := res.SyncedLyrics != ""

	return domain.Lyrics{
		Text:         text,
		IsSynced:     isSynced,
		Instrumental: res.Instrumental,
		Source:       "LRCLIB",
	}, nil
}
