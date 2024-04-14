package config

import "github.com/yannickalex07/dmon/internal/config/destinations"

type DestinationsConfig struct {
	Slack []destinations.SlackDestinationConfig `yaml:"slack"`
}
