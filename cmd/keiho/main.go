package main

import (
	"context"
	"strings"
	"time"

	keiho "github.com/yannickalex07/dmon/pkg"
	"github.com/yannickalex07/dmon/pkg/gcp/dataflow"
	"github.com/yannickalex07/dmon/pkg/local"
	"github.com/yannickalex07/dmon/pkg/slack"
)

func main() {
	ctx := context.Background()

	// build storage
	memoryStorage := local.NewMemoryStorage(time.Hour * 24)

	// build handler
	logHandler := local.LogHandler{}
	slackHandler := slack.SlackHandler{
		Service: slack.NewSlackService("..."),
		Channel: "collection-fna-pipeline-edge-alarms",
	}

	// build checker
	dataflowService := dataflow.NewDataflowService(ctx, "trv-fna-pipeline-edge", "europe-west4", nil)
	dataflowChecker := dataflow.DataflowChecker{Service: dataflowService, Timeout: time.Minute * 2, JobFilter: func(j dataflow.Job) bool {
		return strings.HasPrefix(j.Name, "yannick-")
	}}

	// build monitor
	monitor := keiho.Monitor{
		Storage:  memoryStorage,
		Handlers: []keiho.Handler{&logHandler, &slackHandler},
		Checkers: []keiho.Checker{&dataflowChecker},
	}

	// start monitor
	monitor.StartWithSchedule(ctx, "* * * * *")
}
