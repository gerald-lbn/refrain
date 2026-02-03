package logger

import (
	"log/slog"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLogger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logger Suite")
}

var _ = Describe("Logger", func() {
	Describe("stringToLevel", func() {
		It("should return DebugLevel for 'debug'", func() {
			Expect(stringToLevel("debug")).To(Equal(slog.LevelDebug))
		})

		It("should return InfoLevel for 'info'", func() {
			Expect(stringToLevel("info")).To(Equal(slog.LevelInfo))
		})

		It("should return WarnLevel for 'warn'", func() {
			Expect(stringToLevel("warn")).To(Equal(slog.LevelWarn))
		})

		It("should return ErrorLevel for 'error'", func() {
			Expect(stringToLevel("error")).To(Equal(slog.LevelError))
		})

		It("should return InfoLevel for unknown values", func() {
			Expect(stringToLevel("unknown")).To(Equal(slog.LevelInfo))
			Expect(stringToLevel("")).To(Equal(slog.LevelInfo))
		})
	})
})
