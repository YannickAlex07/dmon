package dataflow_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/dmon/internal/gcp/dataflow"
	gDataflow "google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/option"
)

func TestDataflowServiceListJobs(t *testing.T) {
	// Arrange
	ctx := context.Background()

	udpatedTime := time.Now().UTC().Round(time.Second)
	startTime := udpatedTime.Add(-time.Hour * 1)

	expectedJobs := []dataflow.DataflowJob{
		{
			Id:        "1",
			Name:      "my-job-1",
			Type:      "JOB_TYPE_BATCH",
			StartTime: startTime,
			Status: dataflow.DataflowJobStatus{
				Status:    "JOB_STATE_RUNNING",
				UpdatedAt: udpatedTime,
			},
		},
		{
			Id:        "2",
			Name:      "my-job-2",
			Type:      "JOB_TYPE_STREAMING",
			StartTime: startTime,
			Status: dataflow.DataflowJobStatus{
				Status:    "JOB_STATE_DONE",
				UpdatedAt: udpatedTime,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create response
		resp := &gDataflow.ListJobsResponse{
			Jobs: []*gDataflow.Job{
				{
					Id:               "1",
					Name:             "my-job-1",
					CurrentState:     "JOB_STATE_RUNNING",
					Type:             "JOB_TYPE_BATCH",
					CurrentStateTime: udpatedTime.Format(time.RFC3339),
					StartTime:        startTime.Format(time.RFC3339),
				},
				{
					Id:               "2",
					Name:             "my-job-2",
					CurrentState:     "JOB_STATE_DONE",
					Type:             "JOB_TYPE_STREAMING",
					CurrentStateTime: udpatedTime.Format(time.RFC3339),
					StartTime:        startTime.Format(time.RFC3339),
				},
			},
		}

		// marhsal response
		b, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "unable to marshal request: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Write(b)
	}))

	defer server.Close()

	// create service with overridden endpoint
	service := dataflow.NewDataflowService(ctx, "project", "location", []option.ClientOption{
		option.WithoutAuthentication(),
		option.WithEndpoint(server.URL),
	})

	// Act
	jobs, err := service.ListJobs(ctx)
	if err != nil {
		assert.FailNow(t, "failed to list jobs with: %v", err)
	}

	// Assert
	assert.ElementsMatch(t, expectedJobs, jobs)
}
