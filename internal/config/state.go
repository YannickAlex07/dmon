package config

import (
	"github.com/yannickalex07/inframon/internal/config/states"
	inframon "github.com/yannickalex07/inframon/pkg"
)

type StateConfig struct {
	Memory states.MemoryStateConfig `yaml:"memory"`
}

func (c *StateConfig) Get() (inframon.State, error) {
	return c.Memory.Get(), nil
}
