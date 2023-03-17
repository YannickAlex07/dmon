package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"github.com/yannickalex07/dmon/pkg/api"
	"github.com/yannickalex07/dmon/pkg/config"
	"github.com/yannickalex07/dmon/pkg/handlers/slack"
	"github.com/yannickalex07/dmon/pkg/interfaces"
	"github.com/yannickalex07/dmon/pkg/monitor"
	"github.com/yannickalex07/dmon/pkg/storage"
)

func main() {

	// parse CLI arguments
	configPath := flag.String("c", "./config.yaml", "Path to the config file")
	flag.Parse()

	// parse config
	cfg, err := config.Read(*configPath)
	if err != nil {
		errStr := fmt.Sprintf("Failed to parse config => %s", err.Error())
		panic(errStr)
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

	// build handlers
	handlers := make([]interfaces.Handler, 0)

	// slack handler
	slackHandler := slack.SlackHandler{
		Token:   cfg.Slack.Token,
		Channel: cfg.Slack.Channel,
	}
	handlers = append(handlers, slackHandler)

	// setup state storage
	stateStore := storage.NewMemoryStore(cfg.ExpireTimeoutDuration())

	// setup and start monitor
	monitorFunc := func() {
		monitor.Monitor(*cfg, api.API{}, handlers, stateStore)
	}

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(cfg.RequestInterval).Minute().Do(monitorFunc)
	scheduler.StartBlocking()
}
