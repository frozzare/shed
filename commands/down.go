package commands

import (
	"fmt"

	"github.com/frozzare/shed/docker"
	"github.com/frozzare/shed/log"
	cli "gopkg.in/urfave/cli.v1"
)

var DownCmd = cli.Command{
	Name:   "down",
	Usage:  "down",
	Action: down,
	Flags:  []cli.Flag{},
}

func down(c *cli.Context) {
	app, err := load(c)
	if err != nil {
		log.Error(err)
	}

	fmt.Printf("==>    shed: destroying %s\n", app.Domain())

	// Connect to docker.
	fmt.Println("==>  docker: connecting to docker")
	dock, err := docker.NewDocker(app.Config().Docker)
	if err != nil {
		log.Error(err)
	}

	// Prune removes all unused containers, volumes, networks and images (both dangling and unreferenced).
	fmt.Println("==>  docker: system pruning")
	if err := dock.Prune(); err != nil {
		log.Error(err)
	} else {
		fmt.Println("==>  docker: system pruned")
	}

	// Run docker-compose command.
	cmd := fmt.Sprintf("docker-compose -H %s down -v", dock.Host())
	if err := docker.ExecCmd(cmd, true); err != nil {
		fmt.Printf("==> %s\n", err.Error())
	}

	fmt.Println("==>    shed: done")
}
