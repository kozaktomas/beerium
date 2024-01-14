package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Shops []string `yaml:"shops"`
	Beers []struct {
		Key    string            `yaml:"key"`
		Brand  string            `yaml:"brand"`
		Name   string            `yaml:"name"`
		Type   string            `yaml:"type"`
		Degree int               `yaml:"degree"`
		Size   int               `yaml:"size"`
		Shops  map[string]string `yaml:"shops"`
	} `yaml:"beers"`
}

func LoadConfig(file string) (*Config, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}
	cfg := &Config{}
	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}

	return cfg, nil
}
