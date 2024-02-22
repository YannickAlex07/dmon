package main

import (
	"context"
	"strings"
	"time"

	keiho "github.com/yannickalex07/dmon/pkg"
	checker "github.com/yannickalex07/dmon/pkg/checker"
	dataflow "github.com/yannickalex07/dmon/pkg/external/gcp/dataflow"
	"github.com/yannickalex07/dmon/pkg/external/slack"
	handler "github.com/yannickalex07/dmon/pkg/handler"
	storage "github.com/yannickalex07/dmon/pkg/storage"
)

func main() {
	ctx := context.Background()

	// build storage
	memoryStorage := storage.NewMemoryStorage(time.Hour * 24)

	// build handler
	logHandler := handler.LogHandler{}
	slackHandler := handler.SlackHandler{
		Service: slack.NewSlackService("..."),
		Channel: "collection-fna-pipeline-edge-alarms",
	}

	// build checker
	dataflowService := dataflow.NewDataflowService(ctx, "trv-fna-pipeline-edge", "europe-west4", nil)
	dataflowChecker := checker.DataflowChecker{Service: dataflowService, Timeout: time.Minute * 2, JobFilter: func(j dataflow.Job) bool {
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
