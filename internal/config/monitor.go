package config

type MonitorConfig struct {
	Sources      SourcesConfig      `yaml:"sources"`
	Destinations DestinationsConfig `yaml:"destinations"`
	State        StateConfig        `yaml:"state"`
}
