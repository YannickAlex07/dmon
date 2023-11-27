package checker_test

import (
	"testing"

	"github.com/yannickalex07/dmon/pkg/checker"
)

func TestDataflowChecker(t *testing.T) {
	// Arrange
	_ = checker.DataflowChecker{
		Project:  "my-project",
		Location: "europe-west1",
		JobFilter: func(job checker.DataflowJob) bool {
			return true
		},
	}

	// Act

	// Assert
	t.Fail()
}

func TestDataflowCheckerWithJobFilter(t *testing.T) {
	// Arrange
	_ = checker.DataflowChecker{
		Project:  "my-project",
		Location: "europe-west1",
		JobFilter: func(job checker.DataflowJob) bool {
			return job.Name == "job-1"
		},
	}

	// Act

	// Assert
	t.Fail()
}
