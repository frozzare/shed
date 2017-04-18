package docker

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	api "github.com/fsouza/go-dockerclient"
)

type createContainerOptions struct {
	IP      string
	Image   string
	Name    string
	Ports   []string
	Volumes []string
}

// createOptions will create container options struct
func createOptions(opts *createContainerOptions) api.CreateContainerOptions {
	ip := opts.IP
	if len(ip) == 0 {
		ip = "0.0.0.0"
	}

	publishedPorts := map[api.Port][]api.PortBinding{}
	for _, port := range opts.Ports {
		parts := strings.Split(port, ":")

		if len(parts) < 2 {
			continue
		}

		containerPort := api.Port(parts[1])

		publishedPorts[containerPort+"/tcp"] = []api.PortBinding{api.PortBinding{HostIP: ip, HostPort: parts[0]}}
		publishedPorts[containerPort+"/udp"] = []api.PortBinding{api.PortBinding{HostIP: ip, HostPort: parts[0]}}
	}

	options := api.CreateContainerOptions{
		Name: opts.Name,
		Config: &api.Config{
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
