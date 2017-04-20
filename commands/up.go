package commands

import (
	"fmt"

	"github.com/frozzare/shed/docker"
	"github.com/frozzare/shed/log"
	cli "gopkg.in/urfave/cli.v1"
)

var UpCmd = cli.Command{
	Name:   "up",
	Usage:  "up",
	Action: up,
	Flags:  []cli.Flag{},
}

func up(c *cli.Context) {
	app, err := load(c)
	if err != nil {
		log.Error(err)
	}

	log.Info("shed: creating %s", app.Host())

	// Executing before scripts.
	for _, cmd := range app.Config().BeforeScript {
		if err := docker.ExecCmd(cmd, true); err != nil {
			log.Error(err)
		}
	}

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

	// Start nginx proxy container if it don't exists.
	if err := dock.StartNginxContainer(); err != nil {
		log.Error(err)
	} else {
		log.Info("docker: nginx proxy container is created")
	}

	// Sync application files.
	log.Info("docker: syncing files")
	if err := dock.Sync(); err != nil {
		log.Error(err)
	} else {
		log.Info("docker: syncing done")
	}

	// Run docker-compose commands.
	commands := []string{
		"docker-compose -H %s stop",
		"docker-compose -H %s rm -f",
		"docker-compose -H %s pull",
		"docker-compose -H %s up --build -d",
	}

	for _, cmd := range commands {
		cmd = fmt.Sprintf(cmd, dock.Host())
		if err := docker.ExecCmd(cmd, true); err != nil {
			log.Error(err)
		}
	}

	// Executing after scripts.
	for _, cmd := range app.Config().AfterScript {
		if err := docker.ExecCmd(cmd, true); err != nil {
			log.Error(err)
		}
	}

	log.Info("shed: done, %s is now up", app.URL())
}
