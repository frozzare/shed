package docker

import (
	"fmt"

	api "github.com/fsouza/go-dockerclient"
)

// pullImage will pull a image if it don't exists.
func (d *Docker) pullImage(image string) error {
	if dockerImage, _ := d.client.InspectImage(image); dockerImage == nil {
		fmt.Printf("==> pulling image %s from docker\n", image)
		if err := d.client.PullImage(api.PullImageOptions{
			Repository: image,
			Tag:        "latest",
		}, api.AuthConfiguration{}); err != nil {
			return err
		}
	}

	return nil
}
