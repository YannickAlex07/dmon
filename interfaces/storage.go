package interfaces

import "time"

type Storage interface {
	GetLatestRunTime() time.Time
	SetLatestRunTime(newRunTime time.Time)

	TimeoutAlreadyHandled(id string) bool
	TimeoutHandled(id string)
}
