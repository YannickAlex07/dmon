package sources

import (
	"context"

	inframon "github.com/yannickalex07/inframon/pkg"
	"github.com/yannickalex07/inframon/pkg/gcp/dataflow"
)

type DataflowSourceConfig struct {
	Project  string `validate:"empty=false" yaml:"project"`
	Location string `validate:"empty=false" yaml:"location"`
}

func (c *DataflowSourceConfig) ToSource(ctx context.Context) (inframon.Source, error) {
	service := dataflow.NewDataflowService(ctx, c.Project, c.Location, nil)

	s := dataflow.DataflowSource{
		Service: service,
	}

	return &s, nil
}
