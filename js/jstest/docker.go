//+build !js

package jstest

import (
	"strings"

	"github.com/ory/dockertest/docker"
)

func pullIfNotExists(cli *docker.Client, image string) bool {
	_, err := cli.InspectImage(image)
	if err == nil {
		return true
	}
	i := strings.Index(image, ":")
	err = cli.PullImage(docker.PullImageOptions{
		Repository: image[:i], Tag: image[i+1:],
	}, docker.AuthConfiguration{})
	return err == nil
}
