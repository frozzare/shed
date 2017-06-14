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
		{
			Name:   "down",
			Usage:  "Stop and remove proxy container",
			Action: proxyDown,
		},
		{
			Name:   "up",
			Usage:  "Create proxy container if not created",
			Action: proxyUp,
		},
	},
}

func proxyUp(c *cli.Context) {
	app, err := load(c)
	if err != nil {
		log.Error(err)
	}

	log.Info("shed: creating proxy container")

	// Connect to docker.
	log.Info("docker: connecting to docker")
	dock, err := docker.NewDocker(app.Config().Docker)
	if err != nil {
		log.Error(err)
	}

	// Prune removes all unused containers, volumes, networks and images (both dangling and unreferenced).
	log.Info("docker: system pruning")
	if err := dock.Prune(); err != nil {
		log.Error(err, false)
	} else {
		log.Info("docker: system pruned")
	}

	// Start proxy container if it don't exists.
	if err := dock.StartProxyContainer(); err != nil {
		log.Error(err)
	}

	log.Info("shed: proxy container is up")
}

func proxyDown(c *cli.Context) {
	app, err := load(c)
	if err != nil {
		log.Error(err)
	}

	log.Info("shed: destroying proxy container")

	// Connect to docker.
	log.Info("docker: connecting to docker")
	dock, err := docker.NewDocker(app.Config().Docker)
	if err != nil {
		log.Error(err)
	}

	if err := dock.RemoveContainer("shed_proxy"); err != nil {
		log.Error(err)
	}

	log.Info("shed: proxy container is down")
}
