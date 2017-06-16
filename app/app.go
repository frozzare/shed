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
	os.Setenv("SHED_HOST", a.Host())
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

// Host returns the application host.
func (a *App) Host() string {
	repo := a.Repository()

	for _, name := range []string{repo.Branch, "*"} {
		// Let's allow specific hosts for different branches.
		if len(a.opts.Config.Branches[name].Host) > 0 {
			return a.opts.Config.Branches[name].Host
		}
	}

	// Add leading dot to host name if missing.
	host := a.opts.Config.Host
	if host[0] != '.' {
		host = "." + host
	}

	return repo.Slug + host
}

// URL returns the url to the application.
func (a *App) URL() string {
	scheme := "http"
	port := a.Config().Docker.Proxy.HTTPPort
	if len(a.Config().Docker.Proxy.HTTPSPort) > 0 {
		scheme = "https"
		port = a.Config().Docker.Proxy.HTTPSPort
	}

	if strings.Contains(port, ":") {
		p := strings.Split(port, ":")
		port = p[0]
	}

	if len(port) > 0 {
		port = ":" + port
	}

	if port == "443" || port == "80" {
		port = ""
	}

	hosts := strings.Split(a.Host(), ",")
	for i, host := range hosts {
		hosts[i] = fmt.Sprintf("%s://%s%s", scheme, strings.TrimSpace(host), port)
	}

	return strings.Join(hosts, ", ")
}
