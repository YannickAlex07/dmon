package destinations

import (
	inframon "github.com/yannickalex07/inframon/pkg"
	"github.com/yannickalex07/inframon/pkg/local"
)

type LogDestinationConfig struct {
	Level string `validate:"one_of=debug,info,warn,error" yaml:"level"`
}

func (c *LogDestinationConfig) ToDestination() inframon.Destination {
	d := local.LogDestination{}
	return &d
}
