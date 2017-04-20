package app

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/frozzare/shed/config"
	"github.com/frozzare/shed/repository"
	"github.com/imdario/mergo"
)

var (
	ErrInvalidOptions = errors.New("Invalid options arguments")
)

// App represents the application struct.
type App struct {
	opts *Options
}

// Options represents the application options.
type Options struct {
	Config     config.Config
	Repository repository.Repository
}

// NewApp creates a new application.
func NewApp(opts *Options) (*App, error) {
	if opts == nil {
		return nil, ErrInvalidOptions
	}

	app := &App{opts: opts}

	app.prepare()

	return app, nil
}

// prepare prepares the application.
func (a *App) prepare() {
	os.Setenv("SHED_DOMAIN", a.Domain())
}

// Config returns application config.
func (a *App) Config() config.Config {
	config := a.opts.Config.Branches[a.opts.Repository.Branch]
	mergo.Merge(&config, a.opts.Config)
	return config
}

// Repository returns application repository.
func (a *App) Repository() repository.Repository {
	return a.opts.Repository
}

// Domain returns the application domain.
func (a *App) Domain() string {
	repo := a.Repository()

	// Let's allow specific domains for different branches.
	if len(a.opts.Config.Branches[repo.Branch].Domain) > 0 {
		return a.opts.Config.Branches[repo.Branch].Domain
	}

	// Add leading dot to domain name if missing.
	domain := a.opts.Config.Domain
	if domain[0] != '.' {
		domain = "." + domain
	}

	return repo.Slug + domain
}

// URL returns the url to the application.
func (a *App) URL() string {
	scheme := "http"
	port := a.Config().Docker.Proxy.Ports.HTTP
	if len(a.Config().Docker.Proxy.Ports.HTTPS) > 0 {
		scheme = "https"
		port = a.Config().Docker.Proxy.Ports.HTTPS
	}

	if strings.Contains(port, ":") {
		p := strings.Split(port, ":")
		port = p[0]
	}

	if len(port) > 0 {
		port = ":" + port
	}

	return fmt.Sprintf("%s://%s%s", scheme, a.Domain(), port)
}
