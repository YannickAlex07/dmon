package interfaces

import (
	"github.com/yannickalex07/dmon/models"
)

type Handler interface {
	HandleError(cfg models.Config, job models.Job, messages []models.Message)
	HandleTimeout(cfg models.Config, job models.Job)
}
