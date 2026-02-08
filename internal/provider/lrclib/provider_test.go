package lrclib

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gerald-lbn/refrain/internal/domain"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLrclib(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lrclib Suite")
}

// urlRewriteTransport rewrites the request URL to point to the test server
type urlRewriteTransport struct {
	Target *url.URL
}

func (t *urlRewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = t.Target.Scheme
	req.URL.Host = t.Target.Host
	return http.DefaultTransport.RoundTrip(req)
}

var _ = Describe("Provider", func() {
	var (
		server   *httptest.Server
		provider *Provider
		logger   *slog.Logger
	)

	BeforeEach(func() {
		logger = slog.New(slog.NewTextHandler(GinkgoWriter, nil))
	})

	AfterEach(func() {
		if server != nil {
			server.Close()
		}
	})

	Describe("Search", func() {
		It("should return lyrics when found", func() {
			// Mock Server
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer GinkgoRecover()
				// Verify path is correct (relative to baseURL)
				Expect(r.URL.Path).To(Equal("/api/search"))
				Expect(r.URL.Query().Get("track_name")).To(Equal("Test Song"))

				resp := []lrclibTrack{
					{
						ID:           1,
						TrackName:    "Test Song",
						ArtistName:   "Test Artist",
						PlainLyrics:  "Verse 1",
						SyncedLyrics: "[00:10.00] Verse 1",
					},
				}
				json.NewEncoder(w).Encode(resp)
			}))

			u, err := url.Parse(server.URL)
			Expect(err).NotTo(HaveOccurred())

			client := &http.Client{
				Transport: &urlRewriteTransport{Target: u},
				Timeout:   time.Second,
			}
			provider = New(logger, client)

			results, err := provider.Search(context.Background(), domain.Track{
				Title:  "Test Song",
				Artist: "Test Artist",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))
			Expect(results[0].Text).To(Equal("[00:10.00] Verse 1"))
		})

		It("should handle instrumental tracks correctly", func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp := []lrclibTrack{
					{
						ID:           2,
						TrackName:    "Instrumental Song",
						ArtistName:   "Test Artist",
						Instrumental: true,
					},
				}
				json.NewEncoder(w).Encode(resp)
			}))

			u, _ := url.Parse(server.URL)
			client := &http.Client{Transport: &urlRewriteTransport{Target: u}}
			provider = New(logger, client)

			results, err := provider.Search(context.Background(), domain.Track{Title: "Instrumental Song"})
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))
			Expect(results[0].Text).To(Equal("Instrumental"))
			Expect(results[0].IsSynced).To(BeFalse())
		})
	})

	Describe("Download", func() {
		It("should fetch lyrics by ID", func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer GinkgoRecover()
				Expect(r.URL.Path).To(Equal("/api/get/123"))

				resp := lrclibTrack{
					ID:           123,
					TrackName:    "Test Song",
					PlainLyrics:  "Plain text",
					SyncedLyrics: "",
				}
				json.NewEncoder(w).Encode(resp)
			}))

			u, _ := url.Parse(server.URL)
			client := &http.Client{Transport: &urlRewriteTransport{Target: u}}
			provider = New(logger, client)

			lyrics, err := provider.Download(context.Background(), "123")
			Expect(err).NotTo(HaveOccurred())
			Expect(lyrics.Text).To(Equal("Plain text"))
			Expect(lyrics.IsSynced).To(BeFalse())
		})
	})
})
