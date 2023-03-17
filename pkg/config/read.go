package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func Read(path string) (*Config, error) {
	// open file
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// unmarshal config
	var c *Config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
