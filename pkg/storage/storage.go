package storage

import "time"

type Storage interface {
	GetLatestRuntime() time.Time
	SetLatestRuntime(newRunTime time.Time)

	TimeoutAlreadyHandled(id string) bool
	TimeoutHandled(id string)
}
