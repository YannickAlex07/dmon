package dataflow

import "github.com/yannickalex07/dmon/pkg/model"

type Dataflow interface {
	Jobs() ([]model.Job, error)
	ErrorLogs(jobId string) ([]model.LogEntry, error)
}

type DataflowClient struct {
	Project  string
	Location string
	Prefix   string
}
