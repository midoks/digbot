package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"digbot/cmd"
)

func main() {

	app := cli.NewApp()
	app.Name = "digbot"
	app.Usage = "digbot service"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		cmd.Run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Printf("Failed to start application: %v", err)
	}
}
