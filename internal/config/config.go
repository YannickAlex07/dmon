package config

import (
	"context"

	inframon "github.com/yannickalex07/inframon/pkg"
)

type Config struct {
	Logging LoggingConfig `yaml:"logging"`

	Monitors []MonitorConfig `yaml:"monitors"`
}

func (c *Config) GetMonitors(ctx context.Context) ([]inframon.Monitor, error) {
	var m []inframon.Monitor

	for _, mc := range c.Monitors {
		mon, err := mc.Get(ctx)
		if err != nil {
			return nil, err
		}

		m = append(m, mon)
	}

	return m, nil
}
