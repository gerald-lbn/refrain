package watcher_test

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gerald-lbn/refrain/internal/watcher"
)

var _ = Describe("Watcher", func() {
	var (
		w       *watcher.Watcher
		logger  *slog.Logger
		tempDir string
		ctx     context.Context
		cancel  context.CancelFunc
	)

	BeforeEach(func() {
		logger = slog.New(slog.NewTextHandler(GinkgoWriter, &slog.HandlerOptions{Level: slog.LevelDebug}))

		var err error
		w, err = watcher.New(logger)
		Expect(err).NotTo(HaveOccurred())

		tempDir, err = os.MkdirTemp("", "watcher_test_")
		Expect(err).NotTo(HaveOccurred())

		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		cancel()
		w.Close()
		os.RemoveAll(tempDir)
	})

	Describe("Watch", func() {
		It("should detect new music files", func() {
			ch, err := w.Watch(ctx, []string{tempDir})
			Expect(err).NotTo(HaveOccurred())

			filePath := filepath.Join(tempDir, "song.mp3")

			time.Sleep(50 * time.Millisecond)
			err = os.WriteFile(filePath, []byte("dummy mp3"), 0644)
			Expect(err).NotTo(HaveOccurred())

			Eventually(ch).WithTimeout(2 * time.Second).Should(Receive(Equal(filePath)))
		})

		It("should ignore non-music files", func() {
			ch, err := w.Watch(ctx, []string{tempDir})
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(50 * time.Millisecond)

			txtPath := filepath.Join(tempDir, "readme.txt")
			err = os.WriteFile(txtPath, []byte("hello"), 0644)
			Expect(err).NotTo(HaveOccurred())

			mp3Path := filepath.Join(tempDir, "song.mp3")
			err = os.WriteFile(mp3Path, []byte("dummy mp3"), 0644)
			Expect(err).NotTo(HaveOccurred())

			Eventually(ch).WithTimeout(2 * time.Second).Should(Receive(Equal(mp3Path)))
		})

		It("should auto-watch new subdirectories", func() {
			ch, err := w.Watch(ctx, []string{tempDir})
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(50 * time.Millisecond)

			subDir := filepath.Join(tempDir, "album")
			err = os.Mkdir(subDir, 0755)
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			filePath := filepath.Join(subDir, "track.flac")
			err = os.WriteFile(filePath, []byte("dummy flac"), 0644)
			Expect(err).NotTo(HaveOccurred())

			Eventually(ch).WithTimeout(2 * time.Second).Should(Receive(Equal(filePath)))
		})

		It("should stop when context is cancelled", func() {
			ch, err := w.Watch(ctx, []string{tempDir})
			Expect(err).NotTo(HaveOccurred())

			cancel()

			Eventually(ch).WithTimeout(2 * time.Second).Should(BeClosed())
		})
	})
})
