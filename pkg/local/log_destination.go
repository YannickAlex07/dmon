package local

import (
	"context"
	"log"

	inframon "github.com/yannickalex07/inframon/pkg"
)

type LogDestination struct{}

func (*LogDestination) Handle(ctx context.Context, notification inframon.Notification) error {
	log.Printf("notification: %+v", notification)
	return nil
}
