// Package config loads the service configuration from yaml file.
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ServiceCfg represents the service configuration.
type ServiceCfg struct {
	ServerAddress string `yaml:"server_address"`
}

// Load reads a YAML file and returns a ServiceCfg object.
func Load(path string) (*ServiceCfg, error) {
	data, err := os.ReadFile(path) //nolint:gosec // potential file inclusion
	if err != nil {
		return nil, err
	}
	var cfg ServiceCfg
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
