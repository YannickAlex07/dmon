package dataflow

import (
	"context"

	"github.com/yannickalex07/dmon/pkg/model"
)

type Dataflow interface {
	Jobs(ctx context.Context) ([]model.Job, error)
	ErrorLogs(ctx context.Context, jobId string) ([]model.LogEntry, error)
}

type DataflowClient struct {
	Project  string
	Location string
	Prefix   string
}
