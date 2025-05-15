package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
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
		Engine: DockerEngine,
		client: apiClient,
	}, nil
}

func GetContainerIdByName(dc *DockerController, name string) (string, error) {
	containerList, err := dc.client.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil {
		return "", err
	}
	for _, i := range containerList {
		for _, n := range i.Names {
			if n == name {
				return i.ID, nil
			}
		}
	}
	return "", errors.New("container not found")
}

// exist checks if the image still exist
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

	exist, err := dc.exist(imageName)
	if err != nil {
		return err
	}

	if !exist {
		fmt.Println("Noticed that you don't have the image, pulling now")
		err := dc.update(imageName, name)
		if err != nil {
			return err
		}
	}

	status, err := dc.status("/" + name)
	if err != nil {
		return err
	}

	if status {
		fmt.Println("Noticed there was already a container created, trying to re-starting that container...")
		containerList, err := dc.client.ContainerList(context.Background(), container.ListOptions{
			All: true,
		})
		if err != nil {
			return err
		}

		for _, i := range containerList {
			for _, n := range i.Names {
				if n == "/"+name && i.Image == imageName && i.State != "running" {
					err := dc.client.ContainerStart(context.Background(), i.ID, container.StartOptions{})
					if err != nil {
						return err
					}
					fmt.Println("Container started successfully")
					return nil
				}
			}
		}
		fmt.Println("Noticed container exist, but fail to re-start the container")
	}

	config := &container.Config{
		Image: imageName,
	}

	create, err := dc.client.ContainerCreate(
		context.Background(),
		config,
		&container.HostConfig{PortBindings: nat.PortMap{
			"80/tcp": []nat.PortBinding{{
				HostIP:   "127.0.0.1",
				HostPort: "5000",
			}},
		},
		},
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
	containerList, err := dc.client.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return err
	}

	for _, i := range containerList {
		for _, n := range i.Names {
			if n == "/"+name && i.State == "running" {
				err = dc.client.ContainerStop(context.Background(), n, container.StopOptions{})
				if err != nil {
					return err
				}
				fmt.Println("Container stopped")
				return nil
			}
		}
	}
	return errors.New("container not found")
}

func (dc *DockerController) status(name string) (bool, error) {
	clientList, err := dc.client.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil {
		return false, err
	}
	for _, c := range clientList {
		for _, i := range c.Names {
			if i == name {
				return true, nil
			}
		}
	}
	return false, nil
}

func (dc *DockerController) update(imageName string, name string) error {
	status, err := dc.status("/" + name)
	if err != nil {
		return err
	}
	if status {
		fmt.Println("Old container found")
		id, err := GetContainerIdByName(dc, "/"+name)
		if err != nil {
			return err
		}
		err = dc.client.ContainerRemove(context.Background(), id, container.RemoveOptions{})
		if err != nil {
			return err
		}
		fmt.Println("Container removed")
	}

	exist, err := dc.exist(imageName)
	if err != nil {
		return err
	}
	if exist {
		fmt.Println("Old image found")
		_, err = dc.client.ImageRemove(context.Background(), imageName, image.RemoveOptions{})
		if err != nil {
			return err
		}
		fmt.Println("Image removed")
	}

	pull, err := dc.client.ImagePull(context.Background(), imageName, image.PullOptions{})
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
