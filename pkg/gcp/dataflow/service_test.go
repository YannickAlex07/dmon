package dataflow_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/dmon/pkg/gcp/dataflow"
	gDataflow "google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/option"
)

// LISTING JOBS

func TestDataflowServiceListJobs(t *testing.T) {
	// Arrange
	ctx := context.Background()

	udpatedTime := time.Now().UTC().Round(time.Second)
	startTime := udpatedTime.Add(-time.Hour * 1)

	expectedJobs := []dataflow.Job{
		{
			Id:        "1",
			Name:      "my-job-1",
			Type:      "JOB_TYPE_BATCH",
			Project:   "project",
			Location:  "location",
			StartTime: startTime,
			Status: dataflow.JobStatus{
				Status:    "JOB_STATE_RUNNING",
				UpdatedAt: udpatedTime,
			},
		},
		{
			Id:        "2",
			Name:      "my-job-2",
			Type:      "JOB_TYPE_STREAMING",
			Project:   "project",
			Location:  "location",
			StartTime: startTime,
			Status: dataflow.JobStatus{
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

func TestDataflowServiceListJobsWithInvalidStartTime(t *testing.T) {
	// Arrange
	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create response
		resp := &gDataflow.ListJobsResponse{
			Jobs: []*gDataflow.Job{
				{
					Id:               "1",
					Name:             "my-job-1",
					CurrentState:     "JOB_STATE_RUNNING",
					Type:             "JOB_TYPE_BATCH",
					CurrentStateTime: "invalid-timestamp",
					StartTime:        "invalid-timestamp",
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

	service := dataflow.NewDataflowService(ctx, "project", "location", []option.ClientOption{
		option.WithoutAuthentication(),
		option.WithEndpoint(server.URL),
	})

	// Act
	jobs, err := service.ListJobs(ctx)

	// Assert
	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to parse start time")

	assert.Empty(t, jobs)
}

func TestDataflowServiceListJobsWithInvalidUpdatedTime(t *testing.T) {
	// Arrange
	ctx := context.Background()
	startTime := time.Now().UTC().Round(time.Second)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create response
		resp := &gDataflow.ListJobsResponse{
			Jobs: []*gDataflow.Job{
				{
					Id:               "1",
					Name:             "my-job-1",
					CurrentState:     "JOB_STATE_RUNNING",
					Type:             "JOB_TYPE_BATCH",
					CurrentStateTime: "invalid-timestamp",
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

	service := dataflow.NewDataflowService(ctx, "project", "location", []option.ClientOption{
		option.WithoutAuthentication(),
		option.WithEndpoint(server.URL),
	})

	// Act
	jobs, err := service.ListJobs(ctx)

	// Assert
	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to parse status time")

	assert.Empty(t, jobs)
}

// GETTING ERROR LOGS

func TestDataflowServiceGetLogs(t *testing.T) {
	// Arrange
	ctx := context.Background()
	now := time.Now().UTC().Round(time.Second)

	expectedLogs := []dataflow.LogMessage{
		{
			Text:  "error message",
			Level: dataflow.LEVEL_ERROR,
			Time:  now,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create response
		resp := &gDataflow.ListJobMessagesResponse{
			JobMessages: []*gDataflow.JobMessage{
				{
					MessageImportance: "JOB_MESSAGE_ERROR",
					MessageText:       "error message",
					Time:              now.Format(time.RFC3339),
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

	service := dataflow.NewDataflowService(ctx, "project", "location", []option.ClientOption{
		option.WithoutAuthentication(),
		option.WithEndpoint(server.URL),
	})

	// Act
	logs, err := service.GetLogs(ctx, "my-job", dataflow.LEVEL_ERROR)
	if err != nil {
		assert.FailNow(t, "failed to get error logs with: %v", err)
	}

	// Assert
	assert.ElementsMatch(t, expectedLogs, logs)
}

func TestDataflowServiceGetLogsWithInvalidTime(t *testing.T) {
	// Arrange
	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create response
		resp := &gDataflow.ListJobMessagesResponse{
			JobMessages: []*gDataflow.JobMessage{
				{
					MessageImportance: "JOB_MESSAGE_ERROR",
					MessageText:       "error message",
					Time:              "invalid time",
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

	service := dataflow.NewDataflowService(ctx, "project", "location", []option.ClientOption{
		option.WithoutAuthentication(),
		option.WithEndpoint(server.URL),
	})

	// Act
	logs, err := service.GetLogs(ctx, "my-job", dataflow.LEVEL_ERROR)

	// Assert
	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to parse message time")

	assert.Empty(t, logs)
}
