package main

import (
	"os"

	"github.com/frozzare/shed/shed"
	"github.com/frozzare/shed/version"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "shed"
	app.Version = version.Version
	app.Usage = "command line utility for deploying test containers"
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		shed.DeployCmd,
	}

	app.Run(os.Args)
}
