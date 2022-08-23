package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/yannickalex07/dataflow-monitor/dataflow"
	"github.com/yannickalex07/dataflow-monitor/slack"
)

func handleError(cfg Config, job dataflow.Job) {
	// Get messages from Dataflow API
	// ...

	// Build blocks and send error message
	blocks := slack.ErrorBlocks(job)
	slack.SendMessage(cfg.Slack.Token, cfg.Slack.Channel, blocks)
}

func startMonitor(cfg Config) {
	fmt.Printf("Starting monitoring...\n")

	timeout := 1 * time.Minute // move to config

	// Prepare func
	lastRun := time.Now()
	f := func() {
		// Get All Jobs
		jobs, err := dataflow.ListJobs(cfg.Project.ID, cfg.Project.Location)
		if err != nil {
			return
		}

		// update time
		previousRunTime := lastRun
		lastRun = time.Now()

		// Check Newer Jobs
		for _, job := range jobs {
			// Check Status Updates
			if job.Status.IsNewer(previousRunTime) {
				if job.Status.IsFailed() {
					go handleError(cfg, job)
				}
			}

			// Make sure to cache jobs that are running to long
			if job.Status.IsRunning() {
				runTime := time.Since(job.StartTime)
				if runTime > timeout {
					fmt.Printf("Job %s is running for %s\n", job.Name, runTime)
				}
			}
		}
	}

	// Schedule func
	s := gocron.NewScheduler(time.UTC)
	s.Every(5).Minute().Do(f)
	s.StartBlocking()
}
