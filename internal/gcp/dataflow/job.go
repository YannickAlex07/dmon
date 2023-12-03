package dataflow

import "time"

// Status

// The Status of a Dataflow Job
type DataflowJobStatus struct {
	// The raw status string of the job. You can check this manually
	// or use the provided `IsXXX()`-methods on this struct.
	Status string

	// The time the status was last updated by the Dataflow backend.
	UpdatedAt time.Time
}

// Check if the job has failed according to its status.
func (js DataflowJobStatus) IsFailed() bool {
	return js.Status == "JOB_STATE_FAILED"
}

// Check if the job is running according to its status.
func (js DataflowJobStatus) IsRunning() bool {
	return js.Status == "JOB_STATE_RUNNING"
}

// Job

type DataflowJob struct {
	Id   string
	Name string

	// The raw type of the job. You can check this field manually or use
	// the provided `IsXXX()`-methods on this struct to check it.
	Type string

	// The time that the job started according to the Dataflow backend.
	StartTime time.Time

	// The current status of the Dataflow job, containing the state it is in
	// as well as when it was last updated.
	Status DataflowJobStatus
}

// Check if the job is a streaming job.
func (j DataflowJob) IsStreaming() bool {
	return j.Type == "JOB_TYPE_STREAMING"
}

// Check if the job is a streaming job.
func (j DataflowJob) IsBatch() bool {
	return j.Type == "JOB_TYPE_BATCH"
}

// Check the current runtime of the job. This is calculating by taking the time
// since the start time provided by the Dataflow backend.
func (j DataflowJob) Runtime() time.Duration {
	return time.Since(j.StartTime)
}