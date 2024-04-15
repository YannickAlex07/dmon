package main

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/yannickalex07/inframon/internal/cli"
)

func main() {
	if err := cli.NewApp().Run(os.Args); err != nil {
		pterm.DefaultBasicText.Println(pterm.LightRed("Something went wrong: ", err))
		os.Exit(1)
	}
}
