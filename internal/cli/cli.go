package cli

import (
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
	"github.com/yannickalex07/inframon/internal/config"
)

func run(c *cli.Context) error {
	pterm.DefaultBasicText.Print("Thank you for using Inframon!\n")
	pterm.DefaultBasicText.Println("Please report any issues at github.com/yannickalex07/inframon/issues.")

	// parse config
	configPath := c.String("config")

	pterm.DefaultBasicText.Print("Parsing configuration file...\n")
	config, err := config.Parse(configPath)
	if err != nil {
		return err
	}

	// configure logging
	err = ConfigureLogging(config.Logging)
	if err != nil {
		return err
	}

	// get monitor
	pterm.DefaultBasicText.Print("Creating Monitor...\n")
	mon, err := config.GetMonitor(c.Context)
	if err != nil {
		return err
	}

	// create spinner
	pterm.DefaultBasicText.Println("Starting Monitor...")
	spinner, err := pterm.DefaultSpinner.Start("Running...")
	if err != nil {
		return err
	}

	// starting monitor
	err = nil
	if config.Schedule.Cron != "" {
		err = mon.StartWithSchedule(c.Context, config.Schedule.Cron)
	} else {
		err = mon.Start(c.Context)
	}

	if err != nil {
		spinner.Fail("Monitor failed...\n")
		return err
	}

	spinner.Success("Monitor completed!\n")
	return nil
}

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "inframon"
	app.Usage = "Monitor infrastructure directly from the command line"
	app.Version = "2.0.0"

	app.Commands = []*cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Run the monitoring tool",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "config",
					Aliases:  []string{"c"},
					Usage:    "Path to the configuration file",
					Required: true,
				},
			},
			Action: run,
		},
	}

	return app
}
