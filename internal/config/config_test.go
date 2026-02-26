package config_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gerald-lbn/refrain/internal/config"
)

var _ = Describe("Config", func() {
	AfterEach(func() {
		os.Unsetenv(config.EnvLogLevel)
		os.Unsetenv(config.EnvAppWorkers)
		os.Unsetenv(config.EnvLibraries)
		os.Unsetenv(config.EnvScanInterval)
	})

	It("should use default values when no env vars are set", func() {
		cfg := config.Load()

		Expect(cfg.LogLevel).To(Equal("info"))
		Expect(cfg.Workers).To(Equal(5))
		Expect(cfg.Libraries).To(BeEmpty())
		Expect(cfg.ScanInterval).To(Equal("@every 1h"))
	})

	It("should read log level from env", func() {
		os.Setenv(config.EnvLogLevel, "debug")

		cfg := config.Load()
		Expect(cfg.LogLevel).To(Equal("debug"))
	})

	It("should read workers from env", func() {
		os.Setenv(config.EnvAppWorkers, "10")

		cfg := config.Load()
		Expect(cfg.Workers).To(Equal(10))
	})

	It("should default workers when env value is invalid", func() {
		os.Setenv(config.EnvAppWorkers, "not-a-number")

		cfg := config.Load()
		Expect(cfg.Workers).To(Equal(5))
	})

	It("should default workers when env value is zero or negative", func() {
		os.Setenv(config.EnvAppWorkers, "0")

		cfg := config.Load()
		Expect(cfg.Workers).To(Equal(5))
	})

	It("should parse libraries from comma-separated env", func() {
		os.Setenv(config.EnvLibraries, "/music,/jazz")

		cfg := config.Load()
		Expect(cfg.Libraries).To(HaveLen(2))
		Expect(cfg.Libraries[0].Path).To(Equal("/music"))
		Expect(cfg.Libraries[1].Path).To(Equal("/jazz"))
	})

	It("should trim whitespace from library paths", func() {
		os.Setenv(config.EnvLibraries, " /music , /jazz ")

		cfg := config.Load()
		Expect(cfg.Libraries).To(HaveLen(2))
		Expect(cfg.Libraries[0].Path).To(Equal("/music"))
		Expect(cfg.Libraries[1].Path).To(Equal("/jazz"))
	})

	It("should apply scan interval to all libraries", func() {
		os.Setenv(config.EnvLibraries, "/music")
		os.Setenv(config.EnvScanInterval, "@every 30m")

		cfg := config.Load()
		Expect(cfg.Libraries).To(HaveLen(1))
		Expect(cfg.Libraries[0].ScanInterval).To(Equal("@every 30m"))
	})

	It("should skip empty paths", func() {
		os.Setenv(config.EnvLibraries, "/music,,/jazz,")

		cfg := config.Load()
		Expect(cfg.Libraries).To(HaveLen(2))
	})
})
