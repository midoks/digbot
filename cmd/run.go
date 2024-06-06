package cmd

import (
	"github.com/urfave/cli"
)

var Run = cli.Command{
	Name:        "web",
	Usage:       "This command starts all free proxy web service",
	Description: `Start Free Proxy Service`,
	Action:      WebRun,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "Custom configuration file path"),
	},
}

func WebRun(c *cli.Context) error {

	return nil
}
