package config

import "github.com/yannickalex07/dmon/internal/config/sources"

type SourcesConfig struct {
	HTTP     []sources.HTTPSourceConfig     `yaml:"http"`
	Dataflow []sources.DataflowSourceConfig `yaml:"dataflow"`
}
