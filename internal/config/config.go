package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Routes []Route `yaml:"routes"`
}

type Route struct {
	Path            string           `yaml:"path"`
	Backend         string           `yaml:"backend"`
	RateLimitPolicy *RateLimitPolicy `yaml:"rate_limit_policy"`
}

type RateLimitPolicy struct {
	Requests int           `yaml:"requests"`
	Window   time.Duration `yaml:"window"`
	KeyBy    string        `yaml:"key_by"`
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
