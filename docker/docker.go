package docker

import (
	"os"
	"runtime"
	"strings"

	"github.com/frozzare/shed/config"
	api "github.com/fsouza/go-dockerclient"
)

// Docker represents a docker client.
type Docker struct {
	client *api.Client
}

// NewDocker creates a new docker client.
func NewDocker(config config.Docker) (*Docker, error) {
	var client *api.Client
	var err error

	endpoint := Endpoint(config.Endpoint)

	if len(config.TLSCa) > 0 || len(config.TLSCert) > 0 || len(config.TLSKey) > 0 {
		client, err = api.NewVersionedTLSClient(endpoint, config.TLSCert, config.TLSKey, config.TLSCa, config.Version)
	} else {
		client, err = api.NewVersionedClient(endpoint, config.Version)
	}

	if err != nil {
		return nil, err
	}

	return &Docker{
		client: client,
	}, nil
}

// Endpoint will return the docker endpoint that should be used.
func Endpoint(args ...string) string {
	var endpoint string

	if len(args) > 0 && len(args[0]) > 0 {
		return args[0]
	}

	if os.Getenv("DOCKER_URL") != "" {
		endpoint = os.Getenv("DOCKER_URL")
	} else if runtime.GOOS == "windows" {
		endpoint = "http://localhost:2375"
	} else {
		endpoint = "unix:///var/run/docker.sock"
	}

	return endpoint
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
	image := "jwilder/nginx-proxy"

	// check if image exists or pull it.
	d.pullImage(image)

	// Create container if it don't exists.
	container, err := d.client.CreateContainer(createOptions(&createContainerOptions{
		Name:    "/shed_nginx_proxy",
		Image:   image,
		Ports:   []string{"80:80", "443:433"},
		Volumes: []string{"/var/run/docker.sock:/tmp/docker.sock:ro"},
	}))

	if err != nil {
		// Container already exists error is okey.
		if strings.Contains(err.Error(), "container already exists") {
			return nil
		}

		return err
	}

	return d.startContainer(container)
}
