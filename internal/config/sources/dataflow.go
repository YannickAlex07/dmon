package sources

import (
	"context"
	"time"

	inframon "github.com/yannickalex07/inframon/pkg"
	"github.com/yannickalex07/inframon/pkg/gcp/dataflow"
)

type DataflowSourceConfig struct {
	Project      string `validate:"empty=false" yaml:"project"`
	Location     string `validate:"empty=false" yaml:"location"`
	TimeoutLimit int    `validate:"gt=0" yaml:"timeout_limit"`
}

func (c *DataflowSourceConfig) ToSource(ctx context.Context) (inframon.Source, error) {
	service, err := dataflow.NewDataflowService(ctx, c.Project, c.Location, nil)
	if err != nil {
		return nil, err
	}

	s := dataflow.DataflowSource{
		Service: service,
		Timeout: time.Duration(c.TimeoutLimit) * time.Second,
	}

	return &s, nil
}
