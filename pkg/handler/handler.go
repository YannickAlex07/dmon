package handler

import (
	"context"

	"github.com/yannickalex07/dmon/pkg/model"
)

type Handler interface {
	HandleError(ctx context.Context, job model.Job, entries []model.LogEntry) error
	HandleTimeout(ctx context.Context, job model.Job) error
}
