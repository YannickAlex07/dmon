package cli

import (
	"context"
	"sync"

	"github.com/urfave/cli/v2"
	"github.com/yannickalex07/inframon/internal/config"
	inframon "github.com/yannickalex07/inframon/pkg"
)

func run(c *cli.Context) error {
	// parse config
	config, err := config.Parse("")
	if err != nil {
		return err
	}

	// create context
	ctx := context.Background()

	// get monitors
	mons, err := config.GetMonitors(ctx)
	if err != nil {
		return err
	}

	// run the monitors
	var wg sync.WaitGroup

	for _, mon := range mons {
		wg.Add(1)

		go func(m inframon.Monitor) {
			defer wg.Done()

			m.Start(ctx)
		}(mon)
	}

	go func() {
		wg.Wait()
	}()

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
			Action:  run,
		},
	}

	return app
}
