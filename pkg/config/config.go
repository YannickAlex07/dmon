package config

import "time"

type Config struct {
	RequestInterval int `yaml:"request_interval"`

	Logging struct {
		Verbose bool `yaml:"verbose"`
	}

	Timeout struct {
		MaxTimeout    int `yaml:"max_timeout_duration"`
		ExpireTimeout int `yaml:"expire_timeout_duration"`
	} `yaml:"timeout"`

	Project struct {
		Id       string `yaml:"id"`
		Location string `yaml:"location"`
	} `yaml:"project"`

	Slack struct {
		Token                 string `yaml:"token"`
		Channel               string `yaml:"channel"`
		IncludeErrorSection   bool   `yaml:"include_error_section"`
		IncludeDataflowButton bool   `yaml:"include_dataflow_button"`
	} `yaml:"slack"`
}

func (c Config) MaxTimeoutDuration() time.Duration {
	return time.Duration(c.Timeout.MaxTimeout) * time.Minute
}

func (c Config) ExpireTimeoutDuration() time.Duration {
	return time.Duration(c.Timeout.ExpireTimeout) * time.Minute
}
