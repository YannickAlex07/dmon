package dataflow_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/dmon/pkg/gcp/dataflow"
)

// JOB STATUS

func TestJobStatusIsRunning(t *testing.T) {
	// Arrange
	status := dataflow.JobStatus{
		Status: "JOB_STATE_RUNNING",
	}

	// Assert
	assert.True(t, status.IsRunning())
	assert.False(t, status.IsFailed())
}

func TestJobStatusIsFailed(t *testing.T) {
	// Arrange
	status := dataflow.JobStatus{
		Status: "JOB_STATE_FAILED",
	}

	// Assert
	assert.True(t, status.IsFailed())
	assert.False(t, status.IsRunning())
}

// JOB

func TestJobIsStreaming(t *testing.T) {
	// Arrange
	job := dataflow.Job{
		Id:        "1",
		Name:      "my-job",
		Type:      "JOB_TYPE_STREAMING",
		StartTime: time.Now(),
		Status: dataflow.JobStatus{
			Status:    "JOB_STATE_RUNNING",
			UpdatedAt: time.Now(),
		},
	}

	// Assert
	assert.True(t, job.IsStreaming())
	assert.False(t, job.IsBatch())
}

func TestJobIsBatch(t *testing.T) {
	// Arrange
	job := dataflow.Job{
		Id:        "1",
		Name:      "my-job",
		Type:      "JOB_TYPE_BATCH",
		StartTime: time.Now(),
		Status: dataflow.JobStatus{
			Status:    "JOB_STATE_RUNNING",
			UpdatedAt: time.Now(),
		},
	}

	// Assert
	assert.True(t, job.IsBatch())
	assert.False(t, job.IsStreaming())
}

func TestJobRuntime(t *testing.T) {
	// Arrange
	job := dataflow.Job{
		Id:        "1",
		Name:      "my-job",
		Type:      "JOB_TYPE_BATCH",
		StartTime: time.Now().Add(-time.Hour),
		Status: dataflow.JobStatus{
			Status:    "JOB_STATE_RUNNING",
			UpdatedAt: time.Now(),
		},
	}

	// Act
	runtime := job.Runtime()

	// Assert
	assert.True(t, runtime > time.Hour)
}
