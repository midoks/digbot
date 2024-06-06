package cmd

import (
	"github.com/urfave/cli"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "This command starts web service",
	Description: `Start Web Service`,
	Action:      WebRun,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "Custom configuration file path"),
	},
}

func WebRun(c *cli.Context) error {

	return nil
}
