package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"github.com/yannickalex07/dmon/pkg/config"
	"github.com/yannickalex07/dmon/pkg/dataflow"
	"github.com/yannickalex07/dmon/pkg/handler"
	"github.com/yannickalex07/dmon/pkg/monitor"
	"github.com/yannickalex07/dmon/pkg/storage"
)

func main() {
	ctx := context.Background()

	// parse CLI arguments
	configPath := flag.String("c", "./config.yaml", "Path to the config file")
	flag.Parse()

	// parse config
	cfg, err := config.Read(*configPath)
	if err != nil {
		errStr := fmt.Sprintf("Failed to parse config => %s", err.Error())
		log.Fatal(errStr)
	}

	// setup logging
	if cfg.Logging.Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)

	// build dataflow client
	client := dataflow.DataflowClient{
		Project:  cfg.Project.Id,
		Location: cfg.Project.Location,
	}

	// build handlers
	handlers := make([]handler.Handler, 0)

	// slack handler
	slackHandler := handler.SlackHandler{
		Token:                 cfg.Slack.Token,
		Channel:               cfg.Slack.Channel,
		IncludeErrorSection:   cfg.Slack.IncludeErrorSection,
		IncludeDataflowButton: cfg.Slack.IncludeDataflowButton,
		GCPConfig: handler.SlackGCPConfig{
			Id:       cfg.Project.Id,
			Location: cfg.Project.Location,
		},
	}

	handlers = append(handlers, slackHandler)

	// setup state storage
	stateStore := storage.NewMemoryStore(cfg.ExpireTimeoutDuration())

	// setup and start monitor
	monCfg := monitor.MonitorConfig{
		MaxJobTimeout: cfg.MaxTimeoutDuration(),
	}

	monitorFunc := func() {
		monitor.Monitor(ctx, monCfg, client, handlers, stateStore)
	}

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(cfg.RequestInterval).Minute().Do(monitorFunc)
	scheduler.StartBlocking()
}
