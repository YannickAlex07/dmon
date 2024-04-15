package local

import (
	"context"

	log "github.com/sirupsen/logrus"
	inframon "github.com/yannickalex07/inframon/pkg"
)

type LogDestination struct{}

func (*LogDestination) Handle(ctx context.Context, notification inframon.Notification) error {
	log.Infof("notification: %+v", notification)
	return nil
}
