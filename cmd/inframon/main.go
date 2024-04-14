package main

import (
	"context"
	"strings"
	"time"

	inframon "github.com/yannickalex07/inframon/pkg"
	"github.com/yannickalex07/inframon/pkg/gcp/dataflow"
	"github.com/yannickalex07/inframon/pkg/local"
	"github.com/yannickalex07/inframon/pkg/slack"
)

func main() {
	ctx := context.Background()

	// build storage
	memoryState := local.NewMemoryState(time.Hour * 24)

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
	monitor := inframon.Monitor{
		State:        memoryState,
		Destinations: []inframon.Destination{&logHandler, &slackHandler},
		Sources:      []inframon.Source{&dataflowChecker},
	}

	// start monitor
	monitor.StartWithSchedule(ctx, "* * * * *")
}
