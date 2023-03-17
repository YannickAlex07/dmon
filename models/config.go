package models

import "time"

type Config struct {
	RequestInterval int `default:"10" yaml:"request_interval"`

	Logging struct {
		Verbose bool `default:"false" yaml:"verbose"`
	}

	Timeout struct {
		MaxTimeout    int `default:"60" yaml:"max_timeout_duration"`
		ExpireTimeout int `default:"1440" yaml:"expire_timeout_duration"`
	} `yaml:"timeout"`

	Project struct {
		Id       string `yaml:"id"`
		Location string `yaml:"location"`
	} `yaml:"project"`

	Slack struct {
		Token                 string `yaml:"token"`
		Channel               string `yaml:"channel"`
		IncludeErrorSection   bool   `default:"false" yaml:"include_error_section"`
		IncludeDataflowButton bool   `default:"false" yaml:"include_dataflow_button"`
	} `yaml:"slack"`
}

func (c Config) MaxTimeoutDuration() time.Duration {
	return time.Duration(c.Timeout.MaxTimeout) * time.Minute
}

func (c Config) ExpireTimeoutDuration() time.Duration {
	return time.Duration(c.Timeout.ExpireTimeout) * time.Minute
}
