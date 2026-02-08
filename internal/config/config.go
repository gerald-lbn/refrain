package config

import (
	"github.com/spf13/viper"
)

const DefaultConfigPath = "/config/config.yml"

type Config struct {
	Libraries []LibraryConfig `mapstructure:"libraries"`
	Log       LogConfig       `mapstructure:"log"`
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

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
