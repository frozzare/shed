package docker

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/frozzare/shed/config"
	api "github.com/fsouza/go-dockerclient"
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
		if err := ExecCmd(cmd, false); err != nil {
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
		return errors.New("running on local machine, no need to sync application files")
	}

	cmd := fmt.Sprintf("docker-machine ssh %s -- rm -rf %s", d.config.Machine, os.Getenv("SHED_PATH"))
	if err := ExecCmd(cmd, true); err != nil {
		return err
	}

	cmd = fmt.Sprintf("docker-machine ssh %s -- mkdir -p %s", d.config.Machine, os.Getenv("SHED_PATH"))
	if err := ExecCmd(cmd, true); err != nil {
		return err
	}

	cmd = fmt.Sprintf("docker-machine scp -r . %s:%s", d.config.Machine, os.Getenv("SHED_PATH"))
	if err := ExecCmd(cmd, true); err != nil {
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

// StartNginxContainer will start nginx proxy container.
func (d *Docker) StartNginxContainer() error {
	// Define image.
	image := d.config.Proxy.Image
	if len(image) == 0 {
		image = "jwilder/nginx-proxy"
	}

	// Define ports.
	ports := []string{
		config.Def(d.config.Proxy.Ports.HTTP, "80:80"),
	}

	// Only bind https if https ports is provided.
	if len(d.config.Proxy.Ports.HTTPS) > 0 {
		ports = append(ports, d.config.Proxy.Ports.HTTPS)
	}

	return d.createContainer(&createContainerOptions{
		Name:     "/shed_nginx_proxy",
		Image:    image,
		Recreate: true,
		Ports:    ports,
		Volumes:  []string{"/var/run/docker.sock:/tmp/docker.sock:ro"},
	})
}
