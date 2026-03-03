package config

import (
	"os"
	"strconv"
	"strings"
)

const (
	DefaultAppWorkers = 5
	DefaultLogLevel   = "info"

	EnvAppWorkers = "REFRAIN_APP_WORKERS"
	EnvLogLevel   = "REFRAIN_LOG_LEVEL"
	EnvLibraries  = "REFRAIN_LIBRARIES"
)

type Config struct {
	Libraries []Library
	LogLevel  string
	Workers   int
}

type Library struct {
	Path string
}

// Load reads configuration from environment variables.
func Load() *Config {
	logLevel := getEnv(EnvLogLevel, DefaultLogLevel)

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
					Path: p,
				})
			}
		}
	}

	return &Config{
		Libraries: libraries,
		LogLevel:  logLevel,
		Workers:   workers,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
