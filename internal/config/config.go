package config

import (
	"github.com/spf13/viper"
)

const DefaultConfigPath = "/config/config.yml"
const DefaultAppWorkers = 5

type Config struct {
	Libraries []LibraryConfig `mapstructure:"libraries"`
	Log       LogConfig       `mapstructure:"log"`
	App       AppConfig       `mapstructure:"app"`
}

type AppConfig struct {
	Workers int `mapstructure:"workers"`
}

type LibraryConfig struct {
	Path         string `mapstructure:"path"`
	ScanInterval string `mapstructure:"scan_interval"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

// LoadConfig loads the configuration from the specified file path.
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	viper.SetDefault("app.workers", DefaultAppWorkers)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if config.App.Workers <= 0 {
		config.App.Workers = DefaultAppWorkers
	}

	return &config, nil
}
