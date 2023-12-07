package checker_test

import (
	"context"
	"testing"
)

func TestDataflowChecker(t *testing.T) {
	// Arrange
	_ = context.Background()

	// checker := checker.DataflowChecker{
	// 	Service:   nil,
	// 	JobFilter: func(job dataflow.DataflowJob) bool { return true },
	// 	Timeout:   time.Hour * 1,
	// }

	// // Act
	// checker.Check(ctx, time.Now())

	// Assert
	// t.Fail()
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
