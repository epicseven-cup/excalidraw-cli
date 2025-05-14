package pkg

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// DockerController is the main engine controller that will be used to control the excalidraw container
// Engine -> The value that determine which engine the controller is for
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

func (dc *DockerController) exist(imageName string) (bool, error) {
	containerList, err := dc.client.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, i := range containerList {
		if i.Image == imageName {
			return true, nil
		}
	}
	return false, nil
}

func (dc *DockerController) run(imageName string, name string) error {
	config := &container.Config{
		Image: imageName,
	}
	create, err := dc.client.ContainerCreate(
		context.Background(),
		config,
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

func (dc *DockerController) stop(name string) error {
	err := dc.client.ContainerStop(context.Background(), name, container.StopOptions{})
	if err != nil {
		return err
	}
	return nil
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
