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

type JobsConfig struct {
	ReturnValue []model.Job
	ErrorValue  error
}

type ErrorLogsConfig struct {
	ReturnValue []model.LogEntry
	ErrorValue  error
}

type FakeDataflow struct {
	JobsConfig JobsConfig

	ErrorLogsConfig ErrorLogsConfig
}

func (f FakeDataflow) Jobs(ctx context.Context) ([]model.Job, error) {
	return f.JobsConfig.ReturnValue, f.JobsConfig.ErrorValue
}

func (f FakeDataflow) ErrorLogs(ctx context.Context, jobId string) ([]model.LogEntry, error) {
	return f.ErrorLogsConfig.ReturnValue, f.ErrorLogsConfig.ErrorValue
}

// --- Handler

type HandledErrors struct {
	Job     model.Job
	Entries []model.LogEntry
}

type FakeHandler struct {
	ErrorsHandled   []HandledErrors
	TimeoutsHandled []model.Job
}

func (f *FakeHandler) HandleError(job model.Job, entries []model.LogEntry) error {
	f.ErrorsHandled = append(f.ErrorsHandled, HandledErrors{Job: job, Entries: entries})
	return nil
}

func (f *FakeHandler) HandleTimeout(job model.Job) error {
	f.TimeoutsHandled = append(f.TimeoutsHandled, job)
	return nil
}

// --- StateStore

type ExecutionTimeConfig struct {
	LatestExecutionTimeReturnValue time.Time
	LatestExecutionTimeErrorValue  error

	SetLatestTime time.Time
}

type TimeoutConfig struct {
	WasHandledMap map[string]bool
	Handled       map[string]time.Time
}

type FakeStateStore struct {
	ExecutionTimeConfig ExecutionTimeConfig
	TimeoutConfig       TimeoutConfig
}

func (f FakeStateStore) GetLatestExecutionTime() (time.Time, error) {
	return f.ExecutionTimeConfig.LatestExecutionTimeReturnValue, f.ExecutionTimeConfig.LatestExecutionTimeErrorValue
}
func (f *FakeStateStore) SetLatestExecutionTime(t time.Time) { f.ExecutionTimeConfig.SetLatestTime = t }

func (f FakeStateStore) WasTimeoutHandled(id string) bool     { return f.TimeoutConfig.WasHandledMap[id] }
func (f FakeStateStore) HandleTimeout(id string, t time.Time) { f.TimeoutConfig.Handled[id] = t }

// TESTS

func TestMonitor(t *testing.T) {
	// - Arrange
	lastExecutionTime := time.Now().UTC()

	jobs := []model.Job{
		// Was Updated and Succeeded -> Should trigger nothing
		{
			Id:   "updated-1",
			Name: "updated-1",
			Type: "JOB_TYPE_BATCH",
			Status: model.Status{
				UpdatedAt: lastExecutionTime.Add(1 * time.Minute),
				Status:    "JOB_STATE_DONE",
			},
			StartTime: lastExecutionTime.Add(-1 * time.Hour),
		},
		// Was Updated and Failed -> Should trigger an alert
		{
			Id:   "updated-2",
			Name: "updated-2",
			Type: "JOB_TYPE_BATCH",
			Status: model.Status{
				UpdatedAt: lastExecutionTime.Add(1 * time.Minute),
				Status:    "JOB_STATE_FAILED",
			},
			StartTime: lastExecutionTime.Add(-1 * time.Hour),
		},
		// Was Updated, crossed Timeout and was not handled -> Should trigger an alert
		{
			Id:   "updated-3",
			Name: "updated-3",
			Type: "JOB_TYPE_BATCH",
			Status: model.Status{
				UpdatedAt: lastExecutionTime.Add(1 * time.Minute),
				Status:    "JOB_STATE_RUNNING",
			},
			StartTime: lastExecutionTime.Add(-1 * time.Hour),
		},
		// Was Updated, crossed Timeout but was already handled -> Should trigger nothing
		{
			Id:   "updated-4",
			Name: "updated-4",
			Type: "JOB_TYPE_BATCH",
			Status: model.Status{
				UpdatedAt: lastExecutionTime.Add(1 * time.Minute),
				Status:    "JOB_STATE_RUNNING",
			},
			StartTime: lastExecutionTime.Add(-1 * time.Hour),
		},
		// Crossed Timeout but is streaming -> Should trigger nothing
		{
			Id:   "updated-5",
			Name: "updated-5",
			Type: "JOB_TYPE_STREAMING",
			Status: model.Status{
				UpdatedAt: lastExecutionTime.Add(1 * time.Minute),
				Status:    "JOB_STATE_RUNNING",
			},
			StartTime: lastExecutionTime.Add(-1 * time.Hour),
		},
		// Was not updated -> Should trigger nothing
		{
			Id:   "updated-6",
			Name: "updated-6",
			Type: "JOB_TYPE_BATCH",
			Status: model.Status{
				UpdatedAt: lastExecutionTime.Add(-1 * time.Minute),
				Status:    "JOB_STATE_FAILED",
			},
			StartTime: lastExecutionTime.Add(-1 * time.Hour),
		},
	}

	entries := []model.LogEntry{
		{
			Text: "this is my entry",
			Time: time.Now().UTC(),
		},
	}

	dataflow := FakeDataflow{
		JobsConfig: JobsConfig{
			ReturnValue: jobs,
			ErrorValue:  nil,
		},
		ErrorLogsConfig: ErrorLogsConfig{
			ReturnValue: entries,
			ErrorValue:  nil,
		},
	}

	stateStore := &FakeStateStore{
		ExecutionTimeConfig: ExecutionTimeConfig{
			LatestExecutionTimeReturnValue: lastExecutionTime,
			LatestExecutionTimeErrorValue:  nil,
			SetLatestTime:                  time.Time{},
		},
		TimeoutConfig: TimeoutConfig{
			WasHandledMap: map[string]bool{
				"updated-4": true,
			},
			Handled: map[string]time.Time{},
		},
	}

	fakeHandler := FakeHandler{
		ErrorsHandled:   []HandledErrors{},
		TimeoutsHandled: []model.Job{},
	}

	cfg := monitor.MonitorConfig{
		MaxJobTimeout: 1 * time.Hour,
	}

	// - Act
	err := monitor.Monitor(cfg, dataflow, []handler.Handler{&fakeHandler}, stateStore)

	// - Assert
	assert.Nil(t, err)

	// assert that the latest execution time was updated
	assert.False(t, stateStore.ExecutionTimeConfig.SetLatestTime.IsZero())

	// assert that "updated-3" was put into the state store as handled
	assert.Equal(t, len(stateStore.TimeoutConfig.Handled), 1)

	at, ok := stateStore.TimeoutConfig.Handled["updated-3"]
	assert.True(t, ok)
	assert.False(t, at.IsZero())

	// assert handle error called for "updated-2"
	assert.Equal(t, []HandledErrors{
		{
			Job:     jobs[1],
			Entries: entries,
		},
	}, fakeHandler.ErrorsHandled)

	// assert handle timeout called for "updated-3"
	assert.Equal(t, []model.Job{jobs[2]}, fakeHandler.TimeoutsHandled)
}

func TestMonitorFailsFetchingJobs(t *testing.T) {
	// - Arrange
	dataflow := FakeDataflow{
		JobsConfig: JobsConfig{
			ReturnValue: make([]model.Job, 0),
			ErrorValue:  errors.New("error"),
		},
		ErrorLogsConfig: ErrorLogsConfig{
			ReturnValue: make([]model.LogEntry, 0),
			ErrorValue:  nil,
		},
	}

	stateStore := &FakeStateStore{
		ExecutionTimeConfig: ExecutionTimeConfig{
			LatestExecutionTimeReturnValue: time.Time{},
			LatestExecutionTimeErrorValue:  nil,
			SetLatestTime:                  time.Time{},
		},
		TimeoutConfig: TimeoutConfig{
			WasHandledMap: map[string]bool{},
			Handled:       map[string]time.Time{},
		},
	}

	cfg := monitor.MonitorConfig{
		MaxJobTimeout: 1 * time.Hour,
	}

	// - Act
	err := monitor.Monitor(cfg, dataflow, make([]handler.Handler, 0), stateStore)

	// - Assert
	assert.Error(t, err)
}

func TestMonitorFailingStateStore(t *testing.T) {}

func TestMonitorFailingErrorLogs(t *testing.T) {}
