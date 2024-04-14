package config

type ScheduleConfig struct {
	Cron string `validate:"empty=false" yaml:"cron"`
}
