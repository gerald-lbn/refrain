package config

import (
	"encoding/json"
	"os"
	"strconv"
)

const (
	DefaultConfigPath = "/config/config.json"
	DefaultAppWorkers = 5

	EnvAppWorkers = "REFRAIN_APP_WORKERS"
	EnvLogLevel   = "REFRAIN_LOG_LEVEL"
)

type Config struct {
	Libraries []LibraryConfig `json:"libraries"`
	Log       LogConfig       `json:"log"`
	App       AppConfig       `json:"app"`
}

type AppConfig struct {
	Workers int `json:"workers"`
}

type LibraryConfig struct {
	Path         string `json:"path"`
	ScanInterval string `json:"scan_interval"`
}

type LogConfig struct {
	Level string `json:"level"`
}

// LoadConfig loads the configuration from the specified file path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if workersStr := os.Getenv(EnvAppWorkers); workersStr != "" {
		if workers, err := strconv.Atoi(workersStr); err == nil {
			config.App.Workers = workers
		}
	}

	if logLevel := os.Getenv(EnvLogLevel); logLevel != "" {
		config.Log.Level = logLevel
	}

	if config.App.Workers <= 0 {
		config.App.Workers = DefaultAppWorkers
	}

	return &config, nil
}
