package storage

import "time"

type Storage interface {
	GetLatestExecutionTime() (time.Time, error)
	SetLatestExecutionTime(t time.Time) error

	IsTimeoutStored(id string) (bool, error)
	StoreTimeout(id string, t time.Time) error
}
