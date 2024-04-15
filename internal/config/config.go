package config

import (
	"context"

	inframon "github.com/yannickalex07/inframon/pkg"
)

type Config struct {
	Logging LoggingConfig `yaml:"logging"`

	Sources      SourcesConfig      `yaml:"sources"`
	Destinations DestinationsConfig `yaml:"destinations"`

	State StateConfig `yaml:"state"`

	Schedule ScheduleConfig `yaml:"schedule"`
}

func (c *Config) GetMonitor(ctx context.Context) (inframon.Monitor, error) {
	// get sources
	sources, err := c.Sources.Get(ctx)
	if err != nil {
		return inframon.Monitor{}, err
	}

	// get destinations
	destinations, err := c.Destinations.Get(ctx)
	if err != nil {
		return inframon.Monitor{}, err
	}

	// get state
	state, err := c.State.Get()
	if err != nil {
		return inframon.Monitor{}, err
	}

	mon := inframon.Monitor{
		Sources:      sources,
		Destinations: destinations,
		State:        state,
	}

	return mon, nil
}
