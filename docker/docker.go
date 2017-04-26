package docker

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/frozzare/shed/config"
	"github.com/frozzare/shed/exec"
	api "github.com/fsouza/go-dockerclient"
)

var (
	ErrLocalMachine = errors.New("running on local machine, no need to sync application files")
)

// Docker represents a docker client.
type Docker struct {
	client *api.Client
	config config.Docker
	host   string
}

// NewDocker creates a new docker client.
func NewDocker(config config.Docker) (*Docker, error) {
	var client *api.Client
	var err error
	var host string

	if len(config.Machine) > 0 {
		// Set shed path for docker machine.
		os.Setenv("SHED_PATH", "/tmp/shed")

		// Set docker machine environment variables.
		cmd := fmt.Sprintf("docker-machine env %s", config.Machine)
		if err := exec.Cmd(cmd, false); err != nil {
			return nil, errors.New("docker machine host does not exist: " + config.Machine)
		}

		client, err = api.NewClientFromEnv()
		host = os.Getenv("DOCKER_HOST")
	} else {
		// Set shed path for local machine.
		os.Setenv("SHED_PATH", ".")

		// Find docker host for local machine.
		if os.Getenv("DOCKER_HOST") != "" {
			host = os.Getenv("DOCKER_HOST")
		} else if runtime.GOOS == "windows" {
			host = "http://localhost:2375"
		} else {
			host = "unix:///var/run/docker.sock"
		}

		// Set docker host environment variable so we can use it later.
		os.Setenv("DOCKER_HOST", host)

		client, err = api.NewClient(host)
	}

	if err != nil {
		return nil, err
	}

	return &Docker{
		client: client,
		config: config,
		host:   host,
	}, nil
}

// Sync application files with docker machine.
func (d *Docker) Sync() error {
	if len(d.config.Machine) == 0 {
		return ErrLocalMachine
	}

	cmd := fmt.Sprintf("docker-machine ssh %s -- rm -rf %s", d.config.Machine, os.Getenv("SHED_PATH"))
	if err := exec.Cmd(cmd, true); err != nil {
		return err
	}

	cmd = fmt.Sprintf("docker-machine ssh %s -- mkdir -p %s", d.config.Machine, os.Getenv("SHED_PATH"))
	if err := exec.Cmd(cmd, true); err != nil {
		return err
	}

	cmd = fmt.Sprintf("docker-machine scp -r . %s:%s", d.config.Machine, os.Getenv("SHED_PATH"))
	if err := exec.Cmd(cmd, true); err != nil {
		return err
	}

	return nil
}

// Host will return the docker host that is used.
func (d *Docker) Host() string {
	return d.host
}

// Prune removes all unused containers, volumes, networks and images (both dangling and unreferenced).
func (d *Docker) Prune() error {
	if _, err := d.client.PruneContainers(api.PruneContainersOptions{}); err != nil {
		return err
	}

	if _, err := d.client.PruneImages(api.PruneImagesOptions{}); err != nil {
		return err
	}

	if _, err := d.client.PruneVolumes(api.PruneVolumesOptions{}); err != nil {
		return err
	}

	if _, err := d.client.PruneNetworks(api.PruneNetworksOptions{}); err != nil {
		return err
	}

	return nil
}

// StartProxyContainer will start the proxy container.
func (d *Docker) StartProxyContainer() error {
	// Define image.
	image := d.config.Proxy.Image
	if len(image) == 0 {
		image = "jwilder/nginx-proxy"
	}

	// Define ports.
	ports := []string{
		config.Def(d.config.Proxy.HTTPPort, "80:80"),
	}

	// Only bind https if https ports is provided.
	if len(d.config.Proxy.HTTPSPort) > 0 {
		ports = append(ports, d.config.Proxy.HTTPSPort)
	}

	err := d.createContainer(&createContainerOptions{
		Env:      d.config.Proxy.Env.Values,
		Name:     "/shed_proxy",
		Image:    image,
		Recreate: d.config.Proxy.Recreate,
		Ports:    ports,
		Volumes:  config.DefList(d.config.Proxy.Volumes.Values, []string{"/var/run/docker.sock:/tmp/docker.sock:ro"}),
	})

	if err != nil && strings.Contains(err.Error(), "container already exists") {
		return nil
	}

	return err
}

func (d *Docker) RemoveContainer(name string) error {
	container, err := d.findContainer(name)

	if err != nil {
		return err
	}

	return d.removeContainer(container)
}
