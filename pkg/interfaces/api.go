package interfaces

import "github.com/yannickalex07/dmon/pkg/models"

type API interface {
	Jobs(project string, location string) ([]models.Job, error)
	ErrorLogs(project string, location string, jobId string) ([]models.LogEntry, error)
}
