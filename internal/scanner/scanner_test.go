package scanner_test

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/gerald-lbn/refrain/internal/domain"
	mock_metadata "github.com/gerald-lbn/refrain/internal/mocks/metadata"
	"github.com/gerald-lbn/refrain/internal/scanner"
)

var _ = Describe("Scanner", func() {
	var (
		ctrl       *gomock.Controller
		mockReader *mock_metadata.MockReader
		s          *scanner.Scanner
		logger     *slog.Logger
		tempDir    string
		ctx        context.Context
		cancel     context.CancelFunc
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockReader = mock_metadata.NewMockReader(ctrl)
		logger = slog.New(slog.NewTextHandler(GinkgoWriter, nil))
		s = scanner.New(logger, mockReader)
		ctx, cancel = context.WithCancel(context.Background())

		var err error
		tempDir, err = os.MkdirTemp("", "scanner_test_")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		ctrl.Finish()
		cancel()
		os.RemoveAll(tempDir)
	})

	Describe("Scan", func() {
		It("should find music files and read metadata", func() {
			startData := []byte("dummy mp3 content")
			filePath := filepath.Join(tempDir, "song.mp3")
			err := os.WriteFile(filePath, startData, 0644)
			Expect(err).NotTo(HaveOccurred())

			expectedTrack := &domain.Track{
				Path:     filePath,
				Title:    "Test Song",
				Artist:   "Test Artist",
				Duration: 180 * time.Second,
			}
			mockReader.EXPECT().Read(filePath).Return(expectedTrack, nil).Times(1)

			ch, err := s.Scan(ctx, tempDir)
			Expect(err).NotTo(HaveOccurred())

			var tracks []domain.Track
			for t := range ch {
				tracks = append(tracks, t)
			}

			Expect(tracks).To(HaveLen(1))
			Expect(tracks[0].Path).To(Equal(filePath))
			Expect(tracks[0].Title).To(Equal("Test Song"))
		})

		It("should ignore non-music files", func() {
			filePath := filepath.Join(tempDir, "readme.txt")
			err := os.WriteFile(filePath, []byte("readme"), 0644)
			Expect(err).NotTo(HaveOccurred())

			mockReader.EXPECT().Read(gomock.Any()).Times(0)

			ch, err := s.Scan(ctx, tempDir)
			Expect(err).NotTo(HaveOccurred())

			var tracks []domain.Track
			for t := range ch {
				tracks = append(tracks, t)
			}
			Expect(tracks).To(BeEmpty())
		})
	})
})
