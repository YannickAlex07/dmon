package config

import (
	"os"

	"github.com/creasty/defaults"
	"gopkg.in/dealancer/validate.v2"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logging LoggingConfig `yaml:"logging"`

	Monitors []MonitorConfig `yaml:"monitors"`
}

func Parse(path string) (*Config, error) {
	c := &Config{}

	// set any default values
	if err := defaults.Set(c); err != nil {
		return nil, err
	}

	// read the file
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// unmarshal config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, err
	}

	// validate the config
	if err := validate.Validate(&c); err != nil {
		return nil, err
	}

	return c, nil
}
