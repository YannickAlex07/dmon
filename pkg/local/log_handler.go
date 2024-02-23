package local

import (
	"context"
	"log"

	keiho "github.com/yannickalex07/dmon/pkg"
)

type LogHandler struct{}

func (*LogHandler) Handle(ctx context.Context, notification keiho.Notification) error {
	log.Printf("notification: %+v", notification)
	return nil
}
