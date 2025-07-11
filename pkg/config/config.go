package config

import (
	"os"

	"github.com/caarlos0/env/v9"
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port string `env:"SERVER_PORT" envDefault:"8080"`
}

type LoggerConfig struct {
	Level string `env:"LOGGER_LEVEL" envDefault:"info"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	Logger LoggerConfig `yaml:"logger"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
