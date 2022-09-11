package interfaces

import "github.com/yannickalex07/dmon/models"

type API interface {
	Jobs(project string, location string) ([]models.Job, error)
	Messages(project string, location string, jobId string) ([]models.Message, error)
}
