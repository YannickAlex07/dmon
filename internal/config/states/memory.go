package states

import (
	"time"

	inframon "github.com/yannickalex07/inframon/pkg"
	"github.com/yannickalex07/inframon/pkg/local"
)

type MemoryStateConfig struct {
	TTL int `validate:"gte=1" yaml:"ttl"`
}

func (c *MemoryStateConfig) Get() inframon.State {
	s := local.NewMemoryState(time.Duration(c.TTL))

	return &s
}
