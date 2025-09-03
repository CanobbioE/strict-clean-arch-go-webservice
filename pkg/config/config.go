// Package config loads the service configuration from yaml file.
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the service configuration.
type Config struct {
	ServerAddress string `yaml:"server_address"`
}

// LoadConfig reads a YAML file and returns a Config object.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path) //nolint:gosec // potential file inclusion
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
