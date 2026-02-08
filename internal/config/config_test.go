package config_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gerald-lbn/refrain/internal/config"
)

var _ = Describe("Config", func() {
	var (
		configFile *os.File
		tmpPath    string
		err        error
	)
	BeforeEach(func() {
		configFile, err = os.CreateTemp("", "config-*.yaml")
		Expect(err).NotTo(HaveOccurred())
		tmpPath = configFile.Name()
		content := `
log:
  level: debug
libraries:
  - path: "/tmp/music"
    scan_interval: "1h"
`
		_, err = configFile.WriteString(content)
		Expect(err).NotTo(HaveOccurred())
		err = configFile.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.Remove(tmpPath)
	})

	It("should load configuration from file", func() {
		cfg, err := config.LoadConfig(tmpPath)

		Expect(err).NotTo(HaveOccurred())
		Expect(cfg).NotTo(BeNil())
		Expect(cfg.Log.Level).To(Equal("debug"))
		Expect(cfg.Libraries).To(HaveLen(1))
		Expect(cfg.Libraries[0].Path).To(Equal("/tmp/music"))
		Expect(cfg.Libraries[0].ScanInterval).To(Equal("1h"))
	})

	It("should error when config file does not exist", func() {
		_, err := config.LoadConfig("/non/existent/path/config.yaml")
		Expect(err).To(HaveOccurred())
	})

	It("should error when config file is malformed", func() {
		malformedConfig, err := os.CreateTemp("", "malformed-*.yaml")
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(malformedConfig.Name())

		_, err = malformedConfig.WriteString("libraries: [ unclosed bracket")
		Expect(err).NotTo(HaveOccurred())
		malformedConfig.Close()

		_, err = config.LoadConfig(malformedConfig.Name())
		Expect(err).To(HaveOccurred())
	})
})
