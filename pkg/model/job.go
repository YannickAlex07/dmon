package model

import "time"

// JOB

type Job struct {
	Id        string
	Name      string
	Type      string
	Status    Status
	StartTime time.Time
}

func (j Job) IsStreaming() bool {
	return j.Type == "JOB_TYPE_STREAMING"
}

func (j Job) Runtime() time.Duration {
	return time.Since(j.StartTime)
}

// STATUS

type Status struct {
	Status    string
	UpdatedAt time.Time
}

// func (s Status) IsNewer(t time.Time) bool {
// 	return t.Before(s.UpdatedAt)
// }

func (s Status) IsFailed() bool {
	return s.Status == "JOB_STATE_FAILED"
}

func (s Status) IsCanceled() bool {
	return s.Status == "JOB_STATE_CANCELLED"
}

func (s Status) IsDone() bool {
	return s.Status == "JOB_STATE_DONE"
}

func (s Status) IsRunning() bool {
	return s.Status == "JOB_STATE_RUNNING"
}
