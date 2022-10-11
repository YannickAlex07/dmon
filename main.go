package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/yannickalex07/dmon/api"
	"github.com/yannickalex07/dmon/config"
	"github.com/yannickalex07/dmon/interfaces"
	"github.com/yannickalex07/dmon/monitor"
	"github.com/yannickalex07/dmon/slack"
)

func main() {

	// parse CLI arguments
	configPath := flag.String("config", "./config.yaml", "Path to the config file")
	flag.Parse()

	// parse Config
	cfg, err := config.Read(*configPath)
	if err != nil {
		panic("Failed to parse config with error: error")
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

	// start monitor
	monitor.Start(*cfg, api.API{}, handlers)
}
