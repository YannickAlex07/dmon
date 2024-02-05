package checker_test

import (
	"context"
	"errors"

	dataflow "github.com/yannickalex07/dmon/pkg/gcp/dataflow"
)

// DataflowServiceMock

type DataflowServiceMock struct {
	Jobs []dataflow.Job
	Logs map[string][]dataflow.LogMessage
}

func (s *DataflowServiceMock) ListJobs(ctx context.Context) ([]dataflow.Job, error) {
	return s.Jobs, nil
}

func (s *DataflowServiceMock) GetLogs(ctx context.Context, jobId string, minLevel dataflow.MessageLevel) ([]dataflow.LogMessage, error) {
	logs, ok := s.Logs[jobId]
	if !ok {
		return nil, errors.New("job not found")
	}

	return logs, nil
}
