package helper_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gerald-lbn/refrain/internal/helper"
)

var _ = Describe("Helper", func() {
	Describe("IsMusicFile", func() {
		It("should return true for supported extensions", func() {
			Expect(helper.IsMusicFile("song.mp3")).To(BeTrue())
			Expect(helper.IsMusicFile("song.flac")).To(BeTrue())
			Expect(helper.IsMusicFile("song.m4a")).To(BeTrue())
			Expect(helper.IsMusicFile("song.ogg")).To(BeTrue())
			Expect(helper.IsMusicFile("song.wav")).To(BeTrue())
		})

		It("should return false for unsupported extensions", func() {
			Expect(helper.IsMusicFile("image.jpg")).To(BeFalse())
			Expect(helper.IsMusicFile("document.txt")).To(BeFalse())
			Expect(helper.IsMusicFile("script.sh")).To(BeFalse())
			Expect(helper.IsMusicFile("song")).To(BeFalse())
		})

		It("should be case insensitive", func() {
			Expect(helper.IsMusicFile("song.MP3")).To(BeTrue())
			Expect(helper.IsMusicFile("song.FlAc")).To(BeTrue())
		})
	})
})
