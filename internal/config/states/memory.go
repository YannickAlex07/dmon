package states

type MemoryStateConfig struct {
	TTL int `validate:"gte=1" yaml:"ttl"`
}
