package destinations

import (
	inframon "github.com/yannickalex07/inframon/pkg"
	"github.com/yannickalex07/inframon/pkg/local"
)

type LogDestinationConfig struct {
	Output string `yaml:"output"`
}

func (c *LogDestinationConfig) ToDestination() inframon.Destination {
	d := local.LogDestination{}
	return &d
}
