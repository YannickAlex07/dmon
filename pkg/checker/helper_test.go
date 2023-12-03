package checker_test

import (
	"context"

	dataflow "github.com/yannickalex07/dmon/internal/gcp/dataflow"
)

// DataflowServiceMock

type DataflowServiceMock struct{}

func (*DataflowServiceMock) ListJobs(ctx context.Context) ([]dataflow.DataflowJob, error) {
	return nil, nil
}

func (*DataflowServiceMock) GetLogs(ctx context.Context, jobId string) ([]dataflow.LogEntry, error) {
	return nil, nil
}
