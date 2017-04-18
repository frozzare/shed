package shed

import (
	"fmt"

	"github.com/frozzare/shed/docker"
	cli "gopkg.in/urfave/cli.v1"
)

var DownCmd = cli.Command{
	Name:   "down",
	Usage:  "down",
	Action: down,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "debug",
		},
	},
}

func down(c *cli.Context) {
	app, err := load(c)
	if err != nil {
		rerr(c, err)
		return
	}

	fmt.Printf("==>    shed: destroying %s\n", app.Domain())

	// Run docker-compose command.
	cmd := fmt.Sprintf("docker-compose -H %s down -v", docker.Endpoint())
	if err := docker.ExecCmd(cmd, true); err != nil {
		fmt.Printf("==> %s\n", err.Error())
	}

	fmt.Println("==>    shed: done")
}
