package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/specgen"
	"os"
)

type PodmanController struct {
	Engine ContainerEngine
	conn   context.Context
}

func DeterminePodmanUnixUri(system string) (string, error) {
	// Maybe a switch statement is better here
	if system == "darwin" {
		// macOS
		home := os.Getenv("HOME")
		if home == "" {
			return "", errors.New("HOME environment variable not set")
		}
		return "unix:///" + home + "/.local/share/containers/podman/machine/podman.sock", nil
	} else if system == "linux" {
		return "unix:///var/run/podman/podman.sock", nil
	}
	return "", errors.New("unknown system")
}

func NewPodmanController(system string) (*PodmanController, error) {

	uri, err := DeterminePodmanUnixUri(system)
	if err != nil {
		return nil, err
	}
	conn, err := bindings.NewConnection(context.Background(), uri)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("connected to %s\n", uri)

	return &PodmanController{
		Engine: PodmanEngine,
		conn:   conn,
	}, nil
}

func (pc *PodmanController) run(imageName string, name string) error {
	e, err := pc.exist(name)
	if err != nil {
		return err
	}

	if e {
		fmt.Printf("Noticed that image does not exists yet, using %s\n", name)
		err := pc.update(imageName)
		if err != nil {

			return err
		}
	}

	e, _ = pc.status(name) // It is fine for status to error, since it is a health check so erroring means that it fails to see the container running.
	if e {
		fmt.Println("Noticed that container already exist, removing it")
		err := pc.exit(name)
		if err != nil {
			return err
		}
	}

	s := specgen.NewSpecGenerator(imageName, false)
	s.Name = name
	portMap := types.PortMapping{
		HostIP:        "127.0.0.1",
		ContainerPort: 80,
		HostPort:      5000,
		Range:         1,
		Protocol:      "tcp",
	}
	s.PortMappings = []types.PortMapping{portMap}
	respond, err := containers.CreateWithSpec(pc.conn, s, nil)
	if err != nil {
		return err
	}
	if err := containers.Start(pc.conn, respond.ID, nil); err != nil {
		return err
	}
	fmt.Printf("Container started with\nID: %s\nName:%s\n", respond.ID, name)
	return nil
}

func (pc *PodmanController) exist(imageName string) (bool, error) {
	return images.Exists(pc.conn, imageName, nil)

}

func (pc *PodmanController) exit(name string) error {
	err := containers.Kill(pc.conn, name, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Container killed with Name: %s\n", name)
	rOpt := &containers.RemoveOptions{}
	rOpt = rOpt.WithIgnore(true)
	_, err = containers.Remove(pc.conn, name, rOpt)
	if err != nil {
		return err
	}
	fmt.Println("Container removed")
	return nil
}

func (pc *PodmanController) update(imageName string) error {
	image, err := images.GetImage(pc.conn, imageName, nil)
	if err != nil {
		return err
	}
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
	_, err := containers.RunHealthCheck(pc.conn, name, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}
