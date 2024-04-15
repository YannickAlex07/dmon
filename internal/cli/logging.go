package cli

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/yannickalex07/inframon/internal/config"
)

func ConfigureLogging(config config.LoggingConfig) error {
	// configure level
	log.SetLevel(log.InfoLevel)
	if config.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	// conigure formatter
	log.SetFormatter(&log.TextFormatter{
		DisableColors:    true,
		FullTimestamp:    true,
		DisableQuote:     !config.Quotes,
		DisableTimestamp: !config.Timestamp,
	})

	// configure terminal output
	if config.Terminal {
		log.AddHook(&PTermHook{config.Verbose})
	}

	// configure file output
	if config.File != "" {
		file, err := os.OpenFile(config.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		if err != nil {
			log.Errorf("failed to open log file: %v", err)
			return err
		}

		log.SetOutput(file)
	}

	return nil
}
