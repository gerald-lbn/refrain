package helper_test

import (
	"github.com/gerald-lbn/refrain/internal/helper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Path Helper", func() {
	Describe("ReplaceExtension", func() {
		It("should replace existing extension with new one", func() {
			Expect(helper.ReplaceExtension("song.mp3", ".lrc")).To(Equal("song.lrc"))
			Expect(helper.ReplaceExtension("/path/to/song.flac", ".txt")).To(Equal("/path/to/song.txt"))
		})

		It("should append extension if none exists", func() {
			Expect(helper.ReplaceExtension("song", ".lrc")).To(Equal("song.lrc"))
		})

		It("should handle multiple dots correctly", func() {
			Expect(helper.ReplaceExtension("song.remix.mp3", ".lrc")).To(Equal("song.remix.lrc"))
		})

		It("should return path/ext concatenated if newExt missing dot (though doc says include dot)", func() {
			// Implementation `filepath.Ext` includes dot. `TrimSuffix` removes it.
			// If we pass "lrc" without dot: "song" + "lrc" -> "songlrc"
			Expect(helper.ReplaceExtension("song.mp3", "lrc")).To(Equal("songlrc"))
		})
	})
})
