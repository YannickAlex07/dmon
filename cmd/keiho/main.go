package main

import (
	"context"
	"time"

	keiho "github.com/yannickalex07/dmon/pkg"
	checker "github.com/yannickalex07/dmon/pkg/checker"
	dataflow "github.com/yannickalex07/dmon/pkg/gcp/dataflow"
	handler "github.com/yannickalex07/dmon/pkg/handler"
	storage "github.com/yannickalex07/dmon/pkg/storage"
)

func main() {
	ctx := context.Background()

	// build storage
	memoryStorage := storage.NewMemoryStorage(time.Minute * 5)

	// build handler
	logHandler := handler.LogHandler{}

	// build checker
	dataflowService := dataflow.NewDataflowService(ctx, "trv-master-data-pipeline-edge", "europe-west4", nil)
	dataflowChecker := checker.DataflowChecker{Service: dataflowService, Timeout: time.Hour * 1}

	// build monitor
	monitor := keiho.Monitor{
		Storage:  memoryStorage,
		Handlers: []keiho.Handler{&logHandler},
		Checkers: []keiho.Checker{&dataflowChecker},
	}

	// start monitor
	for {
		monitor.Start(ctx)

		time.Sleep(time.Second * 10)
	}
}
