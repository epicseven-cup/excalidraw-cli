package pkg

import (
	"context"
	"fmt"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerController struct {
	client *client.Client
	Engine ContainerEngine
}

func NewDockerController() (*DockerController, error) {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &DockerController{
		Engine: PodmanEngine,
		client: apiClient,
	}, nil
}

func (dc *DockerController) status(name string) (bool, error) {
	clientList, err := dc.client.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, c := range clientList {
		fmt.Println(c.Names)
		if c.Names[0] == name {
			return true, nil
		}
	}
	return false, nil
}

func (dc *DockerController) exist(name string) (bool, error) {
	imageList, err := dc.client.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, i := range imageList {
		fmt.Println(i)
	}
	return false, nil
}

func (dc *DockerController) run(name string) error {
	dc.client.Container
}
