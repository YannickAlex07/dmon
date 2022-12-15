package interfaces

import (
	"github.com/yannickalex07/dmon/models"
)

type Handler interface {
	HandleError(cfg models.Config, job models.Job, entries []models.LogEntry)
	HandleTimeout(cfg models.Config, job models.Job)
}
