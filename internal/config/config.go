package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Routes []Route `yaml:"routes"`
}

type Route struct {
	Path    string `yaml:"path"`
	Backend string `yaml:"backend"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var config Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
