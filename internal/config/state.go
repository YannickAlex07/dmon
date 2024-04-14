package config

import "github.com/yannickalex07/dmon/internal/config/states"

type StateConfig struct {
	Memory states.MemoryStateConfig `yaml:"memory"`
}
