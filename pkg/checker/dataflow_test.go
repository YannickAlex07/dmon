package checker_test

import (
	"context"
	"testing"
	"time"

	dataflow "github.com/yannickalex07/dmon/internal/gcp/dataflow"
	"github.com/yannickalex07/dmon/pkg/checker"
)

func TestDataflowChecker(t *testing.T) {
	// Arrange
	ctx := context.Background()

	_ = checker.NewDataflowChecker(ctx, "my-project", "europe-west1", func(job dataflow.DataflowJob) bool {
		return true
	}, time.Hour*10)

	// Act

	// Assert
	t.Fail()
}

// func TestDataflowCheckerWithJobFilter(t *testing.T) {
// 	// Arrange
// 	_ = checker.DataflowChecker{
// 		Project:  "my-project",
// 		Location: "europe-west1",
// 		JobFilter: func(job checker.DataflowJob) bool {
// 			return job.Name == "job-1"
// 		},
// 	}

// 	// Act

// 	// Assert
// 	t.Fail()
// }
