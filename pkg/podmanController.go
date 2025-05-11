package pkg

import (
	"context"
	"fmt"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/specgen"
)

type PodmanController struct {
	Engine ContainerEngine
	conn   context.Context
}

func NewPodmanController() (*PodmanController, error) {
	conn, err := bindings.NewConnection(context.Background(), "unix:///run/podman/podman.sock")
	if err != nil {
		return nil, err
	}

	return &PodmanController{
		Engine: PodmanEngine,
		conn:   conn,
	}, nil
}

func (pc *PodmanController) run(name string) error {
	e, err := pc.exist(name)
	if err != nil {
		return err
	}

	if !e {
		fmt.Printf("Noticed that image does not exists yet, using %s", name)
		err := pc.update(name)
		if err != nil {
			return err
		}
	}

	s := specgen.NewSpecGenerator(name, false)
	respond, err := containers.CreateWithSpec(pc.conn, s, nil)
	if err != nil {
		return err
	}
	fmt.Println("Created container with ID:", respond.ID)

	if err := containers.Start(pc.conn, respond.ID, nil); err != nil {
		return err
	}
	fmt.Printf("Container started with ID: %s", respond.ID)
	return nil
}

func (pc *PodmanController) exist(name string) (bool, error) {
	return images.Exists(pc.conn, name, nil)
}

func (pc *PodmanController) exit(name string) error {
	err := containers.Kill(pc.conn, name, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Container killed with ID / Name: %s", name)
	return nil
}

func (pc *PodmanController) update(name string) error {
	image, err := images.GetImage(pc.conn, name, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Image created with Name: %s", name)
	fmt.Printf("ID: %s\n", image.ID)
	fmt.Printf("Architecture: %s\n", image.Architecture)
	fmt.Printf("Annotations: %s\n", image.Annotations)
	fmt.Printf("Author: %s\n", image.Author)
	fmt.Printf("Created: %s\n", image.Created)
	fmt.Printf("Size: %s\n", image.Size)
	fmt.Printf("Version: %s\n", image.Version)
	return nil
}

func (pc *PodmanController) status(name string) (bool, error) {
	return containers.Exists(pc.conn, name, nil)
}
