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
	memoryStorage := local.NewMemoryState(time.Hour * 24)

	// build handler
	logHandler := local.LogDestination{}
	slackHandler := slack.SlackDestination{
		Service: slack.NewSlackService("..."),
		Channel: "collection-fna-pipeline-edge-alarms",
	}

	// build checker
	dataflowService := dataflow.NewDataflowService(ctx, "trv-fna-pipeline-edge", "europe-west4", nil)
	dataflowChecker := dataflow.DataflowSource{Service: dataflowService, Timeout: time.Minute * 2, JobFilter: func(j dataflow.Job) bool {
		return strings.HasPrefix(j.Name, "yannick-")
	}}

	// build monitor
	monitor := keiho.Monitor{
		Storage:      memoryStorage,
		Destinations: []keiho.Handler{&logHandler, &slackHandler},
		Checkers:     []keiho.Checker{&dataflowChecker},
	}

	// start monitor
	monitor.StartWithSchedule(ctx, "* * * * *")
}
