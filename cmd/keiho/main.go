package main

import (
	"context"
	"time"

	dataflow "github.com/yannickalex07/dmon/internal/gcp/dataflow"
	keiho "github.com/yannickalex07/dmon/pkg"
	checker "github.com/yannickalex07/dmon/pkg/checker"
	handler "github.com/yannickalex07/dmon/pkg/handler"
	monitor "github.com/yannickalex07/dmon/pkg/monitor"
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
	monitor := monitor.SingleMonitor{}

	// start monitor
	for {
		monitor.Start(ctx, []keiho.Checker{dataflowChecker}, []keiho.Handler{&logHandler}, memoryStorage)

		time.Sleep(time.Second * 10)
	}
}
