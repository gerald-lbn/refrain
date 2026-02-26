package config

import (
	"os"
	"strconv"
	"strings"
)

const (
	DefaultAppWorkers   = 5
	DefaultLogLevel     = "info"
	DefaultScanInterval = "@every 1h"

	EnvAppWorkers   = "REFRAIN_APP_WORKERS"
	EnvLogLevel     = "REFRAIN_LOG_LEVEL"
	EnvLibraries    = "REFRAIN_LIBRARIES"
	EnvScanInterval = "REFRAIN_SCAN_INTERVAL"
)

type Config struct {
	Libraries    []Library
	LogLevel     string
	Workers      int
	ScanInterval string
}

type Library struct {
	Path         string
	ScanInterval string
}

// Load reads configuration from environment variables.
func Load() *Config {
	logLevel := getEnv(EnvLogLevel, DefaultLogLevel)
	scanInterval := getEnv(EnvScanInterval, DefaultScanInterval)

	workers := DefaultAppWorkers
	if workersStr := os.Getenv(EnvAppWorkers); workersStr != "" {
		if w, err := strconv.Atoi(workersStr); err == nil && w > 0 {
			workers = w
		}
	}

	var libraries []Library
	if paths := os.Getenv(EnvLibraries); paths != "" {
		for p := range strings.SplitSeq(paths, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				libraries = append(libraries, Library{
					Path:         p,
					ScanInterval: scanInterval,
				})
			}
		}
	}

	return &Config{
		Libraries:    libraries,
		LogLevel:     logLevel,
		Workers:      workers,
		ScanInterval: scanInterval,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
