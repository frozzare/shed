package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/frozzare/shed/config"
	api "github.com/fsouza/go-dockerclient"
)

type createContainerOptions struct {
	Env      []string
	Recreate bool
	IP       string
	Image    string
	Name     string
	Ports    []string
	Volumes  []string
}

// createOptions will create container options struct
func createOptions(opts *createContainerOptions) api.CreateContainerOptions {
	ip := config.Def(opts.IP, "0.0.0.0")

	publishedPorts := map[api.Port][]api.PortBinding{}
	for _, port := range opts.Ports {
		parts := strings.Split(port, ":")

		if len(parts) < 2 {
			continue
		}

		containerPort := api.Port(parts[1])

		publishedPorts[containerPort+"/tcp"] = []api.PortBinding{{HostIP: ip, HostPort: parts[0]}}
		publishedPorts[containerPort+"/udp"] = []api.PortBinding{{HostIP: ip, HostPort: parts[0]}}
	}

	options := api.CreateContainerOptions{
		Name: opts.Name,
		Config: &api.Config{
			Env:     opts.Env,
			Image:   opts.Image,
			Volumes: map[string]struct{}{},
		},
		HostConfig: &api.HostConfig{
			Binds:           []string{},
			PublishAllPorts: false,
			PortBindings:    publishedPorts,
			RestartPolicy:   api.AlwaysRestart(),
		},
	}

	for _, volume := range opts.Volumes {
		parts := strings.Split(volume, ":")

		if len(parts) >= 2 {
			if string(parts[0][0]) == "." {
				path, err := os.Getwd()
				if err != nil {
					continue
				}

				volume = filepath.Join(path, volume)
			}

			options.HostConfig.Binds = append(options.HostConfig.Binds, volume)
			options.Config.Volumes[parts[1]] = struct{}{}
		} else {
			options.Config.Volumes[volume] = struct{}{}
		}
	}

	return options
}

// findContainer will find a container with a name.
func (d *Docker) findContainer(name string) (*api.Container, error) {
	containers, err := d.client.ListContainers(api.ListContainersOptions{
		All: true,
	})

	if err != nil {
		return nil, err
	}

	containerName := name
	if containerName[0] != '/' {
		containerName = "/" + containerName
	}

	for _, container := range containers {
		found := false
		for _, name := range container.Names {
			if name == containerName {
				found = true
				break
			}
		}

		if !found {
			continue
		}

		container, err := d.client.InspectContainer(container.ID)
		if err != nil {
			return nil, fmt.Errorf("Failed to inspect container %s: %s", container.ID, err)
		}

		return container, nil
	}

	return nil, nil
}

// createContainer will create a container with the given options.
func (d *Docker) createContainer(opts *createContainerOptions) error {
	// Check if image exists or pull it.
	d.pullImage(opts.Image)

CREATE:
	// Create container if it don't exists.
	container, err := d.client.CreateContainer(createOptions(opts))

	if err != nil {
		// Try to destroy the container if it exists and should be recreated.
		if strings.Contains(err.Error(), "container already exists") && opts.Recreate {
			container, err := d.findContainer(opts.Name)
			if err != nil {
				return err
			}

			if err := d.removeContainer(container); err != nil {
				return err
			}

			goto CREATE
		}

		return err
	}

	return d.startContainer(container)
}

// startContainer will start the container or try to start the container five times before it stops.
func (d *Docker) startContainer(container *api.Container) error {
	attempted := 0
START:
	if err := d.client.StartContainer(container.ID, nil); err != nil {
		// If it is a 500 error it is likely we can retry and be successful.
		if strings.Contains(err.Error(), "API error (500)") {
			if attempted < 5 {
				attempted++
				time.Sleep(1 * time.Second)
				goto START
			}
		}

		return err
	}

	return nil
}

// stopContainer will stop the container or try to stop the container five times before it stops.
func (d *Docker) stopContainer(container *api.Container) error {
	attempted := 0
STOP:
	if err := d.client.StopContainer(container.ID, 0); err != nil {
		if strings.Contains(err.Error(), "API error (500)") {
			if attempted < 5 {
				attempted++
				time.Sleep(1 * time.Second)
				goto STOP
			}
		}

		if strings.Contains(strings.ToLower(err.Error()), "container not running") {
			return nil
		}

		return err
	}

	return nil
}

// removeContainer will remove the container or try to remove the container five times before it stops.
func (d *Docker) removeContainer(container *api.Container) error {
	// No need to remove nil container.
	if container == nil {
		return nil
	}

	attempted := 0
REMOVE:
	if err := d.client.RemoveContainer(api.RemoveContainerOptions{
		ID:    container.ID,
		Force: true,
	}); err != nil {
		if strings.Contains(err.Error(), "API error (500)") {
			if attempted < 5 {
				attempted++
				time.Sleep(1 * time.Second)
				goto REMOVE
			}
		}

		return err
	}

	return nil
}
