package local

import (
	"context"

	log "github.com/sirupsen/logrus"
	inframon "github.com/yannickalex07/inframon/pkg"
)

type LogDestination struct {
	Level log.Level
}

func NewLogDestination(level string) (*LogDestination, error) {
	l, err := log.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	return &LogDestination{
		Level: l,
	}, nil
}

func (d *LogDestination) Handle(ctx context.Context, notification inframon.Notification) error {
	log.StandardLogger().Logf(d.Level, "Received Notification: %+v", notification)
	return nil
}
