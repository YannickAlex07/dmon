package gmon

import "context"

type Handler interface {
	Handle(ctx context.Context, notification Notification) error
}
