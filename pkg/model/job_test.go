package model_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/dmon/pkg/model"
)

func TestJobIsStreaming(t *testing.T) {
	// - Arrange
	job := model.Job{
		Id:   "my-job-id",
		Name: "my-job",
		Type: "JOB_TYPE_STREAMING",
		Status: model.Status{
			Status:    "STARTED",
			UpdatedAt: time.Now(),
		},
		StartTime: time.Now().Add(-1 * time.Hour),
	}

	// - Assert
	assert.True(t, job.IsStreaming())
}

func TestJobRuntime(t *testing.T) {
	// - Arrange
	job := model.Job{
		Id:   "my-job-id",
		Name: "my-job",
		Type: "JOB_TYPE_STREAMING",
		Status: model.Status{
			Status:    "JOB_STATE_RUNNING",
			UpdatedAt: time.Now(),
		},
		StartTime: time.Now().Add(-2 * time.Hour),
	}

	// - Act
	runtime := job.Runtime()

	// - Assert

	// testing this is very tricky because Runtime() will call
	// time.Since(startTime) and this of course includes the runtime
	// of the test between creating the job and calling .Runtime() on it.
	// Because that time is no known to us here, we will just check that the
	// Runtime is between our expected time and +1 minute of that time.
	lowerBound := 2 * time.Hour
	upperBound := 2*time.Hour + 1*time.Minute

	assert.True(t, runtime > lowerBound)
	assert.True(t, runtime < upperBound)
}

func TestStatusIsFailed(t *testing.T) {
	// - Arrange
	status := model.Status{
		Status:    "JOB_STATE_FAILED",
		UpdatedAt: time.Now(),
	}

	// - Assert
	assert.True(t, status.IsFailed())
}

func TestStatusIsCanceled(t *testing.T) {
	// - Arrange
	status := model.Status{
		Status:    "JOB_STATE_CANCELLED",
		UpdatedAt: time.Now(),
	}

	// - Assert
	assert.True(t, status.IsCanceled())
}

func TestStatusIsDone(t *testing.T) {
	// - Arrange
	status := model.Status{
		Status:    "JOB_STATE_DONE",
		UpdatedAt: time.Now(),
	}

	// - Assert
	assert.True(t, status.IsDone())
}

func TestStatusIsRunning(t *testing.T) {
	// - Arrange
	status := model.Status{
		Status:    "JOB_STATE_RUNNING",
		UpdatedAt: time.Now(),
	}

	// - Assert
	assert.True(t, status.IsRunning())
}
