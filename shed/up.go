package shed

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
	fmt.Println("==>  docker: connecting to docker remote api")
	dock, err := docker.NewDocker(app.Config().Docker)
	if err != nil {
		rerr(c, err)
		return
	}

	// Prune removes all unused containers, volumes, networks and images (both dangling and unreferenced).
	fmt.Println("==>  docker: system pruning")
	if err := dock.Prune(); err != nil {
		fmt.Printf("==>    error: %s\n", err.Error())
	} else {
		fmt.Println("==>  docker: system pruned")
	}

	// Start nginx proxy container if it don't exists.
	if err := dock.StartNginxContainer(); err != nil {
		fmt.Printf("==>    error: %s\n", err.Error())
	} else {
		fmt.Println("==>  docker: nginx proxy container is created or already exists")
	}

	// Run docker-compose command.
	cmd := fmt.Sprintf("docker-compose -H %s up -d --force-recreate", docker.Endpoint())
	if err := docker.ExecCmd(cmd, true); err != nil {
		fmt.Printf("==> %s\n", err.Error())
	}

	fmt.Printf("==>    shed: done, http://%s is now up\n", app.Domain())
}
