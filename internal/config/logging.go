package config

type LoggingConfig struct {
	Terminal bool   `default:"false" yaml:"terminal"`
	File     string `yaml:"file"`

	Verbose bool `default:"false" yaml:"verbose"`

	Timestamp bool `default:"true" yaml:"timestamp"`
	Quotes    bool `default:"true" yaml:"quotes"`
}
