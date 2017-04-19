package commands

import (
	"fmt"

	"github.com/frozzare/shed/docker"
	cli "gopkg.in/urfave/cli.v1"
)

var UpCmd = cli.Command{
	Name:   "up",
	Usage:  "up",
	Action: up,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "debug",
		},
	},
}

func up(c *cli.Context) {
	app, err := load(c)
	if err != nil {
		rerr(c, err)
		return
	}

	fmt.Printf("==>    shed: creating %s\n", app.Domain())

	// Connect to docker.
	fmt.Println("==>  docker: connecting to docker")
	dock, err := docker.NewDocker(app.Config().Docker)
	if err != nil {
		rerr(c, err)
		return
	}

	// Prune removes all unused containers, volumes, networks and images (both dangling and unreferenced).
	fmt.Println("==>  docker: system pruning")
	if err := dock.Prune(); err != nil {
		fmt.Printf("==>   error: %s\n", err.Error())
	} else {
		fmt.Println("==>  docker: system pruned")
	}

	// Start nginx proxy container if it don't exists.
	if err := dock.StartNginxContainer(); err != nil {
		fmt.Printf("==>   error: %s\n", err.Error())
	} else {
		fmt.Println("==>  docker: nginx proxy container is created")
	}

	// Sync application files.
	fmt.Println("==>  docker: syncing files")
	if err := dock.Sync(); err != nil {
		fmt.Printf("==>   error: %s\n", err.Error())
	} else {
		fmt.Println("==>  docker: syncing done")
	}

	// Run docker-compose command.
	flags := "--force-recreate"
	if app.Config().Docker.Build {
		flags = "--build"
	}
	cmd := fmt.Sprintf("docker-compose -H %s up -d %s", dock.Host(), flags)
	if err := docker.ExecCmd(cmd, true); err != nil {
		fmt.Printf("==> %s\n", err.Error())
	}

	fmt.Printf("==>    shed: done, %s is now up\n", app.URL())
}
