package pkg

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type DockerController struct {
	client *client.Client
	Engine ContainerEngine
}

func NewDockerController(system string) (*DockerController, error) {
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
	create, err := dc.client.ContainerCreate(
		context.Background(),
		&container.Config{Image: name},
		nil,
		nil,
		nil,
		name,
	)
	if err != nil {
		return err
	}
	fmt.Printf("Container created: %s\n", create.ID)
	err = dc.client.ContainerStart(context.Background(), create.ID, container.StartOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Container with ID: %s started\n", create.ID)
	return nil
}

func (dc *DockerController) update(name string) error {
	pull, err := dc.client.ImagePull(context.Background(), name, image.PullOptions{})
	if err != nil {
		return err
	}
	fmt.Println("Image pulled")
	err = pull.Close()
	if err != nil {
		return err
	}
	return nil
}

func (dc *DockerController) exit(name string) error {
	err := dc.client.ContainerRemove(context.Background(), name, container.RemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}
