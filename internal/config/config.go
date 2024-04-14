package config

type Config struct {
	Logging LoggingConfig `yaml:"logging"`

	Sources      SourcesConfig      `yaml:"sources"`
	Destinations DestinationsConfig `yaml:"destinations"`
	State        StateConfig        `yaml:"state"`
}
