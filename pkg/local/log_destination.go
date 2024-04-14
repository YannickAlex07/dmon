package local

import (
	"context"
	"log"

	keiho "github.com/yannickalex07/dmon/pkg"
)

type LogDestination struct{}

func (*LogDestination) Handle(ctx context.Context, notification keiho.Notification) error {
	log.Printf("notification: %+v", notification)
	return nil
}
