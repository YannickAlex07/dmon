package monitor_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/dmon/pkg/handler"
	"github.com/yannickalex07/dmon/pkg/model"
	"github.com/yannickalex07/dmon/pkg/monitor"
)

// FAKES

// --- Dataflow

type FakeJob struct {
	Job     model.Job
	Entries []model.LogEntry
}

type FakeDataflow struct {
	FakeJobs []FakeJob

	JobsFetchError    error
	EntriesFetchError error
}

func (f FakeDataflow) Jobs(ctx context.Context) ([]model.Job, error) {
	jobs := []model.Job{}

	for _, j := range f.FakeJobs {
		jobs = append(jobs, j.Job)
	}

	return jobs, f.JobsFetchError
}

func (f FakeDataflow) ErrorLogs(ctx context.Context, jobId string) ([]model.LogEntry, error) {
	for _, j := range f.FakeJobs {
		if j.Job.Id == jobId {
			return j.Entries, f.EntriesFetchError
		}
	}

	return []model.LogEntry{}, f.EntriesFetchError
}

// --- Handler

type HandledErrors struct {
	Job     model.Job
	Entries []model.LogEntry
}

type FakeHandler struct {
	HandledErrors   []HandledErrors
	HandledTimeouts []model.Job

	HandleErrorError   error
	HandleTimeoutError error
}

func (f *FakeHandler) HandleError(ctx context.Context, job model.Job, entries []model.LogEntry) error {
	f.HandledErrors = append(f.HandledErrors, HandledErrors{Job: job, Entries: entries})
	return f.HandleErrorError
}

func (f *FakeHandler) HandleTimeout(ctx context.Context, job model.Job) error {
	f.HandledTimeouts = append(f.HandledTimeouts, job)
	return f.HandleTimeoutError
}

// --- StateStore

type ExecutionTimeConfig struct {
	GetValue time.Time
	GetError error

	SetValue time.Time
	SetError error
}

type TimeoutConfig struct {
	IsStoredMap   map[string]bool
	IsStoredError error

	Stored      map[string]time.Time
	StoredError error
}

type FakeStateStore struct {
	ExecutionTimeConfig ExecutionTimeConfig
	TimeoutConfig       TimeoutConfig
}

func (f FakeStateStore) GetLatestExecutionTime(ctx context.Context) (time.Time, error) {
	return f.ExecutionTimeConfig.GetValue, f.ExecutionTimeConfig.GetError
}
func (f *FakeStateStore) SetLatestExecutionTime(ctx context.Context, t time.Time) error {
	f.ExecutionTimeConfig.SetValue = t
	return f.ExecutionTimeConfig.SetError
}

func (f FakeStateStore) IsTimeoutStored(ctx context.Context, id string) (bool, error) {
	return f.TimeoutConfig.IsStoredMap[id], f.TimeoutConfig.IsStoredError
}
func (f FakeStateStore) StoreTimeout(ctx context.Context, id string, t time.Time) error {
	f.TimeoutConfig.Stored[id] = t
	return f.TimeoutConfig.StoredError
}

// TESTS

// This test will assert the general logic of the monitor when no component fails.
// This means that neither the state store, handler or dataflow client fail with any of their
// requests. Therefore this tests is basically asserting all the core logic.
func TestMonitor(t *testing.T) {
	// - Arrange
	ctx := context.Background()
	lastExecutionTime := time.Now().UTC()

	jobs := []FakeJob{
		// Was Updated and Succeeded -> Should trigger nothing
		{
			Job: model.Job{

				Id:   "updated-1",
				Name: "updated-1",
				Type: "JOB_TYPE_BATCH",
				Status: model.Status{
					UpdatedAt: lastExecutionTime.Add(1 * time.Minute),
					Status:    "JOB_STATE_DONE",
				},
				StartTime: lastExecutionTime.Add(-1 * time.Hour),
			},
			Entries: []model.LogEntry{},
		},
		// Was Updated and Failed -> Should trigger an alert
		{
			Job: model.Job{
				Id:   "updated-2",
				Name: "updated-2",
				Type: "JOB_TYPE_BATCH",
				Status: model.Status{
					UpdatedAt: lastExecutionTime.Add(1 * time.Minute),
					Status:    "JOB_STATE_FAILED",
				},
				StartTime: lastExecutionTime.Add(-1 * time.Hour),
			},
			Entries: []model.LogEntry{},
		},
		// Was Updated, crossed Timeout and was not handled -> Should trigger an alert
		{
			Job: model.Job{
				Id:   "updated-3",
				Name: "updated-3",
				Type: "JOB_TYPE_BATCH",
				Status: model.Status{
					UpdatedAt: lastExecutionTime.Add(1 * time.Minute),
					Status:    "JOB_STATE_RUNNING",
				},
				StartTime: lastExecutionTime.Add(-1 * time.Hour),
			},
			Entries: []model.LogEntry{},
		},
		// Was Updated, crossed Timeout but was already handled -> Should trigger nothing
		{
			Job: model.Job{

				Id:   "updated-4",
				Name: "updated-4",
				Type: "JOB_TYPE_BATCH",
				Status: model.Status{
					UpdatedAt: lastExecutionTime.Add(1 * time.Minute),
					Status:    "JOB_STATE_RUNNING",
				},
				StartTime: lastExecutionTime.Add(-1 * time.Hour),
			},
			Entries: []model.LogEntry{},
		},
		// Crossed Timeout but is streaming -> Should trigger nothing
		{
			Job: model.Job{

				Id:   "updated-5",
				Name: "updated-5",
				Type: "JOB_TYPE_STREAMING",
				Status: model.Status{
					UpdatedAt: lastExecutionTime.Add(1 * time.Minute),
					Status:    "JOB_STATE_RUNNING",
				},
				StartTime: lastExecutionTime.Add(-1 * time.Hour),
			},
			Entries: []model.LogEntry{},
		},
		// Was not updated -> Should trigger nothing
		{
			Job: model.Job{

				Id:   "updated-6",
				Name: "updated-6",
				Type: "JOB_TYPE_BATCH",
				Status: model.Status{
					UpdatedAt: lastExecutionTime.Add(-1 * time.Minute),
					Status:    "JOB_STATE_FAILED",
				},
				StartTime: lastExecutionTime.Add(-1 * time.Hour),
			},
			Entries: []model.LogEntry{},
		},
	}

	dataflow := FakeDataflow{
		FakeJobs: jobs,

		JobsFetchError:    nil,
		EntriesFetchError: nil,
	}

	stateStore := &FakeStateStore{
		ExecutionTimeConfig: ExecutionTimeConfig{
			GetValue: lastExecutionTime,
			GetError: nil,
			SetValue: time.Time{},
			SetError: nil,
		},
		TimeoutConfig: TimeoutConfig{
			IsStoredMap: map[string]bool{
				"updated-4": true,
			},
			IsStoredError: nil,
			Stored:        map[string]time.Time{},
			StoredError:   nil,
		},
	}

	fakeHandler := FakeHandler{
		HandledErrors:      []HandledErrors{},
		HandleErrorError:   nil,
		HandledTimeouts:    []model.Job{},
		HandleTimeoutError: nil,
	}

	cfg := monitor.MonitorConfig{
		MaxJobTimeout: 1 * time.Hour,
	}

	// - Act
	err := monitor.Monitor(ctx, cfg, dataflow, []handler.Handler{&fakeHandler}, stateStore)

	// - Assert
	assert.Nil(t, err)

	// assert that the latest execution time was updated
	assert.False(t, stateStore.ExecutionTimeConfig.SetValue.IsZero())

	// assert that "updated-3" was put into the state store as handled
	assert.Equal(t, len(stateStore.TimeoutConfig.Stored), 1)

	at, ok := stateStore.TimeoutConfig.Stored["updated-3"]
	assert.True(t, ok)
	assert.False(t, at.IsZero())

	// assert handle error called for "updated-2"
	assert.Equal(t, []HandledErrors{
		{
			Job:     jobs[1].Job,
			Entries: jobs[1].Entries,
		},
	}, fakeHandler.HandledErrors)

	// assert handle timeout called for "updated-3"
	assert.Equal(t, []model.Job{jobs[2].Job}, fakeHandler.HandledTimeouts)
}

// This tests asserts the monitor behavior when the client fails to fetch
// any jobs from the Dataflow API.
// Expected is that the monitor just returns an error.
func TestMonitorFailsToFetchJobs(t *testing.T) {
	// - Arrange
	ctx := context.Background()

	dataflow := FakeDataflow{
		FakeJobs:          []FakeJob{},
		JobsFetchError:    errors.New("error"),
		EntriesFetchError: nil,
	}

	stateStore := &FakeStateStore{
		ExecutionTimeConfig: ExecutionTimeConfig{
			GetValue: time.Time{},
			GetError: nil,
			SetValue: time.Time{},
			SetError: nil,
		},
		TimeoutConfig: TimeoutConfig{
			IsStoredMap:   map[string]bool{},
			IsStoredError: nil,
			Stored:        map[string]time.Time{},
			StoredError:   nil,
		},
	}

	cfg := monitor.MonitorConfig{
		MaxJobTimeout: 1 * time.Hour,
	}

	// - Act
	err := monitor.Monitor(ctx, cfg, dataflow, make([]handler.Handler, 0), stateStore)

	// - Assert
	assert.Error(t, err)
}

// This test asserts what happens when we fail to fetch the error logs for a given job.
// Expected is that the handled will just receive empty error logs.
func TestMonitorFailsToFetchErrorLogs(t *testing.T) {
	// - Arrange
	ctx := context.Background()

	lastExecutionTime := time.Now().UTC()

	jobs := []FakeJob{
		{
			Job: model.Job{

				Id:   "job",
				Name: "job",
				Type: "JOB_TYPE_BATCH",
				Status: model.Status{
					UpdatedAt: lastExecutionTime.Add(1 * time.Minute),
					Status:    "JOB_STATE_FAILED",
				},
				StartTime: lastExecutionTime,
			},
			Entries: []model.LogEntry{
				{
					Text: "This is a failed log entry",
					Time: time.Now(),
				},
			},
		},
	}

	dataflow := FakeDataflow{
		FakeJobs: jobs,

		JobsFetchError:    nil,
		EntriesFetchError: errors.New("error"),
	}

	stateStore := &FakeStateStore{
		ExecutionTimeConfig: ExecutionTimeConfig{
			GetValue: lastExecutionTime,
			GetError: nil,
			SetValue: time.Time{},
			SetError: nil,
		},
		TimeoutConfig: TimeoutConfig{
			IsStoredMap:   map[string]bool{},
			IsStoredError: nil,
			Stored:        map[string]time.Time{},
			StoredError:   nil,
		},
	}

	fakeHandler := FakeHandler{
		HandledErrors:   []HandledErrors{},
		HandledTimeouts: []model.Job{},
	}

	cfg := monitor.MonitorConfig{
		MaxJobTimeout: 1 * time.Hour,
	}

	// - Act
	err := monitor.Monitor(ctx, cfg, dataflow, []handler.Handler{&fakeHandler}, stateStore)

	// - Assert
	assert.Nil(t, err)

	assert.Equal(t, []HandledErrors{
		{
			Job:     jobs[0].Job,
			Entries: []model.LogEntry{}, // -> Should be empty
		},
	}, fakeHandler.HandledErrors)
}

// This test asserts what happens when the monitor fails to fetch
// the last execution from the state store.
// Expected is simply that we will take the current time as the last execution time.
func TestMonitorFailsToFetchExecutionTime(t *testing.T) {
	// - Arrange
	ctx := context.Background()

	lastExecutionTime := time.Now().UTC().Add(-1 * time.Hour)

	jobs := []FakeJob{
		// This job failed after the last execution time
		// However we will let the store fail to return the last execution time
		// which will trigger the monitor to use the current time. This should
		// result into this error not being reported.
		{
			Job: model.Job{

				Id:   "job",
				Name: "job",
				Type: "JOB_TYPE_BATCH",
				Status: model.Status{
					UpdatedAt: lastExecutionTime.Add(1 * time.Minute),
					Status:    "JOB_STATE_FAILED",
				},
				StartTime: lastExecutionTime,
			},
			Entries: []model.LogEntry{
				{
					Text: "This is a failed log entry",
					Time: time.Now(),
				},
			},
		},
	}

	dataflow := FakeDataflow{
		FakeJobs: jobs,

		JobsFetchError:    nil,
		EntriesFetchError: nil,
	}

	stateStore := &FakeStateStore{
		ExecutionTimeConfig: ExecutionTimeConfig{
			GetValue: lastExecutionTime,
			GetError: errors.New("error"), // this should trigger monitor to use now()
			SetValue: time.Time{},
			SetError: nil,
		},
		TimeoutConfig: TimeoutConfig{
			IsStoredMap:   map[string]bool{},
			IsStoredError: nil,
			Stored:        map[string]time.Time{},
			StoredError:   nil,
		},
	}

	fakeHandler := FakeHandler{
		HandledErrors:   []HandledErrors{},
		HandledTimeouts: []model.Job{},
	}

	cfg := monitor.MonitorConfig{
		MaxJobTimeout: 10 * time.Hour,
	}

	// - Act
	err := monitor.Monitor(ctx, cfg, dataflow, []handler.Handler{&fakeHandler}, stateStore)

	// - Assert
	assert.Nil(t, err)
	assert.Equal(t, []HandledErrors{}, fakeHandler.HandledErrors) // -> We should see no failed jobs
}

// This test asserts what happens when the monitor fails to check if a given timeout
// is already stored within the state store.
// Expected is that we just assume it was never handled.
func TestMonitorFailsTimeoutStoredCheck(t *testing.T) {
	// - Arrange
	ctx := context.Background()

	lastExecutionTime := time.Now().UTC().Add(-1 * time.Hour)

	jobs := []FakeJob{
		// This job runs for over an hour already, which should prompt the monitor
		// to issue a timeout alert.
		{
			Job: model.Job{

				Id:   "job",
				Name: "job",
				Type: "JOB_TYPE_BATCH",
				Status: model.Status{
					UpdatedAt: lastExecutionTime,
					Status:    "JOB_STATE_RUNNING",
				},
				StartTime: lastExecutionTime,
			},
			Entries: []model.LogEntry{},
		},
	}

	dataflow := FakeDataflow{
		FakeJobs: jobs,

		JobsFetchError:    nil,
		EntriesFetchError: nil,
	}

	stateStore := &FakeStateStore{
		ExecutionTimeConfig: ExecutionTimeConfig{
			GetValue: lastExecutionTime,
			GetError: nil,
			SetValue: time.Time{},
			SetError: nil,
		},
		TimeoutConfig: TimeoutConfig{
			IsStoredMap: map[string]bool{
				// -> The timeout was already stored in the past
				// meaning we should not send an alert - HOWEVER! we will
				// let the check for this value fail with an error, which
				// should trigger a new alert
				"job": true,
			},
			IsStoredError: errors.New("error"), // -> this should trigger the monitor to issue a new alert
			Stored:        map[string]time.Time{},
			StoredError:   nil,
		},
	}

	fakeHandler := FakeHandler{
		HandledErrors:   []HandledErrors{},
		HandledTimeouts: []model.Job{},
	}

	cfg := monitor.MonitorConfig{
		MaxJobTimeout: 1 * time.Minute,
	}

	// - Act
	err := monitor.Monitor(ctx, cfg, dataflow, []handler.Handler{&fakeHandler}, stateStore)

	// - Assert
	assert.Nil(t, err)
	assert.Equal(t, []model.Job{jobs[0].Job}, fakeHandler.HandledTimeouts)
}
