package handler

import "github.com/yannickalex07/dmon/pkg/model"

type Handler interface {
	HandleError(job model.Job, entries []model.LogEntry)
	HandleTimeout(job model.Job)
}
