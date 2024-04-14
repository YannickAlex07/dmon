package main

import (
	"os"

	"github.com/yannickalex07/inframon/internal/cli"
)

func main() {
	if err := cli.NewApp().Run(os.Args); err != nil {
		os.Exit(1)
	}
}
