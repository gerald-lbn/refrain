package logger_test

import (
	"context"
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gerald-lbn/refrain/internal/logger"
)

var _ = Describe("Logger", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("New", func() {
		DescribeTable("should optionally set the correct log level",
			func(levelStr string, expectedLevel slog.Level) {
				log := logger.New(levelStr)
				Expect(log).NotTo(BeNil())
				Expect(log.Enabled(ctx, expectedLevel)).To(BeTrue(), "Expected level %v to be enabled", expectedLevel)

				var expectedDisabledLevel slog.Level
				switch expectedLevel {
				case slog.LevelInfo:
					expectedDisabledLevel = slog.LevelDebug
				case slog.LevelWarn:
					expectedDisabledLevel = slog.LevelInfo
				case slog.LevelError:
					expectedDisabledLevel = slog.LevelWarn
				}

				if expectedLevel > slog.LevelDebug {
					Expect(log.Enabled(ctx, expectedDisabledLevel)).To(BeFalse(), "Expected level %v to be disabled", expectedDisabledLevel)
				}
			},
			Entry("debug level", "debug", slog.LevelDebug),
			Entry("info level", "info", slog.LevelInfo),
			Entry("warn level", "warn", slog.LevelWarn),
			Entry("error level", "error", slog.LevelError),
			Entry("upper case debug", "DEBUG", slog.LevelDebug),
			Entry("mixed case info", "InFo", slog.LevelInfo),
			Entry("invalid level defaults to info", "unknown", slog.LevelInfo),
			Entry("empty level defaults to info", "", slog.LevelInfo),
		)
	})
})
