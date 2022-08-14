package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/yannickalex07/dataflow-monitor/dataflow"
)

func main() {
	projectId := "trv-master-data-pipeline-prod"
	location := "europe-west4"
	// timeout := 3 * time.Hour

	lastRun := time.Now()
	f := func() {
		// Get All Jobs
		jobs, err := dataflow.ListJobs(projectId, location)
		if err != nil {
			return
		}

		// Check Newer Jobs
		for _, job := range jobs {
			// Check Status Updates
			if job.Status.IsNewer(lastRun) {
				println("Job " + job.Name + " has newer status")
			}

			// Check Run Length
			runLength := time.Since(job.StartTime)

			fmt.Printf("Job %s is running for %s\n", job.Name, runLength)
			// if runLength > timeout {
			// 	println("Job " + job.Name + " ran for too long")
			// }
		}

		// update time
		lastRun = time.Now()
	}

	s := gocron.NewScheduler(time.UTC)
	s.Every(30).Seconds().Do(f)
	s.StartBlocking()
}
