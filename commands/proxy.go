package commands

import (
	"github.com/frozzare/shed/docker"
	"github.com/frozzare/shed/log"
	cli "gopkg.in/urfave/cli.v1"
)

var ProxyCmd = cli.Command{
	Name:  "proxy",
	Usage: "",
	Subcommands: []cli.Command{
		cli.Command{
			Name:   "down",
			Usage:  "Stop and remove proxy container",
			Action: proxyDown,
		},
	},
}

func proxyDown(c *cli.Context) {
	app, err := load(c)
	if err != nil {
		log.Error(err)
	}

	log.Info("shed: destroying proxy container")

	dock, err := docker.NewDocker(app.Config().Docker)
	if err != nil {
		log.Error(err)
	}

	if err := dock.RemoveContainer("shed_proxy"); err != nil {
		log.Error(err)
	}

	log.Info("shed: proxy container is destroyed")
}
