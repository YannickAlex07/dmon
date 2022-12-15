package interfaces

import "time"

type Storage interface {
	GetLatestRunTime() time.Time
	SetLatestRunTime(newRunTime time.Time)

	WasTimeoutHandled(id string) bool
	TimeoutHandled(id string)
}
