package commands

import (
	"github.com/frozzare/shed/config"
	"github.com/frozzare/shed/docker"
	"github.com/frozzare/shed/exec"
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
	cfg := app.Config()

	log.Info("shed: creating %s", app.Host())

	// Executing before scripts.
	exec.CmdList(cfg.BeforeScript, true)

	// Connect to docker.
	log.Info("docker: connecting to docker")
	dock, err := docker.NewDocker(cfg.Docker)
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
	} else {
		log.Info("docker: proxy container is created")
	}

	// Sync application files.
	log.Info("docker: syncing files")
	if err := dock.Sync(); err != nil {
		if err == docker.ErrLocalMachine {
			log.Error(err, false)
		} else {
			log.Error(err)
		}
	} else {
		log.Info("docker: syncing done")
	}

	// Run docker-compose commands.
	commands := config.DefList(cfg.Script, []string{
		"docker-compose -H %s stop",
		"docker-compose -H %s rm -f",
		"docker-compose -H %s pull",
		"docker-compose -H %s up --build -d",
	})
	exec.CmdList(commands, true)

	// Executing after scripts.
	exec.CmdList(cfg.AfterScript, true)

	log.Info("shed: done, %s is now up", app.URL())
}
