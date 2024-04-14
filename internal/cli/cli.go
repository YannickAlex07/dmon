package cli

import "github.com/urfave/cli/v2"

func run(c *cli.Context) error {
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
