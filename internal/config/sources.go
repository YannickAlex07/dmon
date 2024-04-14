package config

import (
	"context"
	"errors"

	"github.com/yannickalex07/inframon/internal/config/sources"
	inframon "github.com/yannickalex07/inframon/pkg"
)

type SourcesConfig struct {
	HTTP     []sources.HTTPSourceConfig     `yaml:"http"`
	Dataflow []sources.DataflowSourceConfig `yaml:"dataflow"`
}

func (c *SourcesConfig) Get(ctx context.Context) ([]inframon.Source, error) {
	var s []inframon.Source

	// for _, h := range c.HTTP {
	// 	s = append(s, h.ToSource())
	// }

	for _, d := range c.Dataflow {
		ds, err := d.ToSource(ctx)
		if err != nil {
			return nil, err
		}

		s = append(s, ds)
	}

	if len(s) == 0 {
		return nil, errors.New("no sources configured")
	}

	return s, nil
}
