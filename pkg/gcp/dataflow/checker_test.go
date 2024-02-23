package dataflow_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	keiho "github.com/yannickalex07/dmon/pkg"
	dataflow "github.com/yannickalex07/dmon/pkg/gcp/dataflow"
)

func TestDataflowChecker(t *testing.T) {
	// Arrange
	ctx := context.Background()
	since := time.Now().UTC().Add(-time.Hour * 1)

	service := &DataflowServiceMock{
		Jobs: []dataflow.Job{
			// This job should not trigger an alert
			// as updatedAt before the last runtime (since)
			{
				Id:        "a",
				Name:      "job-1",
				Project:   "project",
				Location:  "location",
				StartTime: since.Add(-time.Hour * 1),
				Status: dataflow.JobStatus{
					Status:    "JOB_STATE_FAILED",
					UpdatedAt: since.Add(-time.Minute * 1),
				},
			},
			// This job should not trigger an alert
			// as it is not failed
			{
				Id:        "b",
				Name:      "job-2",
				Project:   "project",
				Location:  "location",
				StartTime: since.Add(-time.Hour * 1),
				Status: dataflow.JobStatus{
					Status:    "JOB_STATE_RUNNING",
					UpdatedAt: since.Add(time.Minute * 30),
				},
			},
			// This job should trigger an alert
			// as it failed after our last runtime (since)
			{
				Id:        "c",
				Name:      "job-3",
				Project:   "project",
				Location:  "location",
				StartTime: since.Add(-time.Hour * 1),
				Status: dataflow.JobStatus{
					Status:    "JOB_STATE_FAILED",
					UpdatedAt: since.Add(time.Minute * 1),
				},
			},
			// This job should trigger an alert
			// as it is running longer than our timeout limit
			{
				Id:        "d",
				Name:      "job-4",
				Project:   "project",
				Location:  "location",
				StartTime: since.Add(-time.Hour * 5),
				Status: dataflow.JobStatus{
					Status:    "JOB_STATE_RUNNING",
					UpdatedAt: since.Add(-time.Hour * 5),
				},
			},
		},
		Logs: map[string][]dataflow.LogMessage{
			"c": {
				{
					Text:  "This is a log message",
					Level: dataflow.LEVEL_ERROR,
					Time:  time.Now(),
				},
			},
		},
	}

	// parse the error URL for later checks
	cUrl, err := url.Parse("https://console.cloud.google.com/dataflow/jobs/location/c?project=project&authuser=1&hl=en")
	if err != nil {
		t.Fatal(err)
	}

	dUrl, err := url.Parse("https://console.cloud.google.com/dataflow/jobs/location/d?project=project&authuser=1&hl=en")
	if err != nil {
		t.Fatal(err)
	}

	expectedNotifications := []keiho.Notification{
		// This is the notification for job id "c"
		{
			Key:         fmt.Sprintf("DATAFLOW-ERROR-c-%s", since.Add(-time.Hour*1).Format(time.RFC3339)),
			Title:       "❌ Dataflow Job Failed",
			Description: fmt.Sprintf("The job `job-3` with id `c` failed at *%s*!", since.Add(time.Minute*1).Format(time.RFC1123)),
			Logs: []string{
				"This is a log message",
			},
			Links: map[string]*url.URL{
				"Open In Dataflow": cUrl,
			},
		},
		// This is the notification for job id "d"
		{
			Key:         fmt.Sprintf("DATAFLOW-TIMEOUT-d-%s", since.Add(-time.Hour*5).Format(time.RFC3339)),
			Title:       "⏱️ Dataflow Job Running For Too Long",
			Description: "The job `job-4` with id `d` crossed the maximum timeout limit with a runtime of *6h0m0s*.",
			Logs:        []string{},
			Links: map[string]*url.URL{
				"Open In Dataflow": dUrl,
			},
		},
	}

	checker := dataflow.DataflowChecker{
		Service:   service,
		JobFilter: func(job dataflow.Job) bool { return true },
		Timeout:   time.Hour * 3,
	}

	// Act
	notifications, err := checker.Check(ctx, since)

	// Assert
	assert.NoError(t, err)

	assert.Equal(t, expectedNotifications, notifications)
}

func TestDataflowCheckerWithJobFilter(t *testing.T) {
	// Arrange
	ctx := context.Background()
	since := time.Now().UTC().Add(-time.Hour * 1)

	service := &DataflowServiceMock{
		Jobs: []dataflow.Job{
			// This job should trigger an alert
			// as it failed after our last runtime (since)
			{
				Id:        "a",
				Name:      "job-1",
				Project:   "project",
				Location:  "location",
				StartTime: since.Add(-time.Hour * 1),
				Status: dataflow.JobStatus{
					Status:    "JOB_STATE_FAILED",
					UpdatedAt: since.Add(time.Minute * 1),
				},
			},
			// This job should **not** trigger an alert
			// as it should be ignored by our job filter
			{
				Id:        "b",
				Name:      "job-2",
				Project:   "project",
				Location:  "location",
				StartTime: since.Add(-time.Hour * 1),
				Status: dataflow.JobStatus{
					Status:    "JOB_STATE_FAILED",
					UpdatedAt: since.Add(time.Minute * 1),
				},
			},
		},
		Logs: map[string][]dataflow.LogMessage{
			"a": {
				{
					Text:  "This is a log message",
					Level: dataflow.LEVEL_ERROR,
					Time:  time.Now(),
				},
			},
		},
	}

	uiUrl, err := url.Parse("https://console.cloud.google.com/dataflow/jobs/location/a?project=project&authuser=1&hl=en")
	if err != nil {
		t.Fatal(err)
	}

	expectedNotifications := []keiho.Notification{
		// This is the notification for job id "c"
		{
			Key:         fmt.Sprintf("DATAFLOW-ERROR-a-%s", since.Add(-time.Hour*1).Format(time.RFC3339)),
			Title:       "❌ Dataflow Job Failed",
			Description: fmt.Sprintf("The job `job-1` with id `a` failed at *%s*!", since.Add(time.Minute*1).Format(time.RFC1123)),
			Logs: []string{
				"This is a log message",
			},
			Links: map[string]*url.URL{
				"Open In Dataflow": uiUrl,
			},
		},
	}

	filter := func(job dataflow.Job) bool {
		// filter every job that is named "job-2"
		return job.Name != "job-2"
	}

	checker := dataflow.DataflowChecker{
		Service:   service,
		JobFilter: filter,
		Timeout:   time.Hour * 3,
	}

	// Act
	notifications, err := checker.Check(ctx, since)

	// Assert
	assert.NoError(t, err)

	assert.Equal(t, expectedNotifications, notifications)
}
