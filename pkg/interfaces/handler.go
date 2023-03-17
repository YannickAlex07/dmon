package interfaces

import (
	"github.com/yannickalex07/dmon/pkg/models"
)

type Handler interface {
	HandleError(job models.Job, entries []models.LogEntry)
	HandleTimeout(job models.Job)
}
