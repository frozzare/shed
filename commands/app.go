package commands

import (
	"fmt"

	"github.com/frozzare/shed/app"
	"github.com/frozzare/shed/config"
	"github.com/frozzare/shed/log"
	"github.com/frozzare/shed/repository"
	cli "gopkg.in/urfave/cli.v1"
)

var AppCmd = cli.Command{
	Name:   "app",
	Usage:  "Shows information about the shed application",
	Action: appAction,
	Flags:  []cli.Flag{},
}

func load(c *cli.Context) (*app.App, error) {
	// Load configuration.
	config, err := config.NewConfig(c.GlobalString("file"))
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

func appAction(c *cli.Context) {
	app, err := load(c)
	if err != nil {
		log.Error(err)
	}

	fmt.Printf("Branch: %s\n", app.Repository().Branch)
	fmt.Printf("URL: %s\n", app.URL())
}
