package interfaces

import "github.com/yannickalex07/dmon/pkg/models"

type DataflowClient interface {
	Jobs() ([]models.Job, error)
	ErrorLogs(jobId string) ([]models.LogEntry, error)
}
