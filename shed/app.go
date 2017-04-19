package shed

import (
	"fmt"

	"github.com/frozzare/shed/app"
	"github.com/frozzare/shed/config"
	"github.com/frozzare/shed/repository"
	cli "gopkg.in/urfave/cli.v1"
)

var AppCmd = cli.Command{
	Name:   "app",
	Usage:  "app",
	Action: ap,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "debug",
		},
	},
}

func load(c *cli.Context) (*app.App, error) {
	// Load configuration.
	config, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	// Load git repository.
	repo, err := repository.NewRepository(config.Git)
	if err != nil {
		return nil, err
	}

	// Create a new application.
	return app.NewApp(&app.Options{
		Config:     config,
		Repository: repo,
	})
}

func ap(c *cli.Context) {
	app, err := load(c)
	if err != nil {
		rerr(c, err)
		return
	}

	fmt.Printf("Branch: %s\n", app.Repository().Branch)
	fmt.Printf("URL: %s\n", app.URL())
}
