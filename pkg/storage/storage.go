package storage

import "time"

type Storage interface {
	GetLatestExecutionTime() (time.Time, error)
	SetLatestExecutionTime(t time.Time)

	WasTimeoutHandled(id string) bool
	HandleTimeout(id string, t time.Time)
}
