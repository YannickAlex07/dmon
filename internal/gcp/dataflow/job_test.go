package dataflow_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/dmon/internal/gcp/dataflow"
)

func TestJobIsStreaming(t *testing.T) {
	// Arrange
	job := dataflow.DataflowJob{
		Id:        "1",
		Name:      "my-job",
		Type:      "JOB_TYPE_STREAMING",
		StartTime: time.Now(),
		Status: dataflow.DataflowJobStatus{
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
	job := dataflow.DataflowJob{
		Id:        "1",
		Name:      "my-job",
		Type:      "JOB_TYPE_BATCH",
		StartTime: time.Now(),
		Status: dataflow.DataflowJobStatus{
			Status:    "JOB_STATE_RUNNING",
			UpdatedAt: time.Now(),
		},
	}

	// Assert
	assert.True(t, job.IsBatch())
	assert.False(t, job.IsStreaming())
}
