package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/yannickalex07/dataflow-monitor/dataflow"
)

func startMonitor(cfg Config) {
	timeout := 1 * time.Minute // move to config

	// Prepare func
	lastRun := time.Now()
	f := func() {
		// Get All Jobs
		jobs, err := dataflow.ListJobs(cfg.Project.ID, cfg.Project.Location)
		if err != nil {
			return
		}

		// Check Newer Jobs
		for _, job := range jobs {
			// Check Status Updates
			if job.Status.IsNewer(lastRun) {
				if job.Status.IsFailed() {
					fmt.Printf("Job %s failed\n", job.Name)
				} else if job.Status.IsDone() {
					fmt.Printf("Job %s finished\n", job.Name)
				}
			}

			if job.Status.IsRunning() {
				runTime := time.Since(job.StartTime)
				if runTime > timeout {
					fmt.Printf("Job %s is running for %s\n", job.Name, runTime)
				}
			}
		}

		// update time
		lastRun = time.Now()
	}

	// Schedule func
	s := gocron.NewScheduler(time.UTC)
	s.Every(30).Seconds().Do(f)
	s.StartBlocking()
}
