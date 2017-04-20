package main

import (
	"os"

	"github.com/frozzare/shed/commands"
	"github.com/frozzare/shed/version"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "shed"
	app.Version = version.Version
	app.Usage = "cli for deploying test containers based on a git repository"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Specify an alternate shed file (default: .shed.yml or shed.yml)",
		},
	}
	app.Commands = []cli.Command{
		commands.AppCmd,
		commands.DownCmd,
		commands.UpCmd,
	}

	app.Run(os.Args)
}
