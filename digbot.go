package main

import (
	"log"
	"os"

	"digbot/cmd"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "digbot"
	app.Usage = "digbot service"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		cmd.Web,
		cmd.Scan,
	}

	if err := app.Run(os.Args); err != nil {
		log.Printf("Failed to start application: %v", err)
	}
}
