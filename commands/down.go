package commands

import (
	"fmt"

	"github.com/frozzare/shed/docker"
	"github.com/frozzare/shed/exec"
	"github.com/frozzare/shed/log"
	cli "gopkg.in/urfave/cli.v1"
)

var DownCmd = cli.Command{
	Name:   "down",
	Usage:  "Stop and remove containers, networks, images, and volumes",
	Action: downAction,
	Flags:  []cli.Flag{},
}

func downAction(c *cli.Context) {
	app, err := load(c)
	if err != nil {
		log.Error(err)
	}

	log.Info("shed: destroying %s", app.Host())

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

	// Run docker-compose command.
	cmd := fmt.Sprintf("docker-compose -H %s down -v", dock.Host())
	if err := exec.Cmd(cmd, true); err != nil {
		log.Error(err)
	}

	log.Info("shed: done")
}
