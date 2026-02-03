package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Library struct {
	Path string `yaml:"path"`
	Name string `yaml:"name"`
}

type Logger struct {
	Level string `yaml:"level"`
}

type Config struct {
	Libraries []Library `yaml:"libraries"`
	Logger    Logger    `yaml:"logger"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
