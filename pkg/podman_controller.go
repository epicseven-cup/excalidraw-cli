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
	return &PodmanController{
		Engine: PodmanEngine,
		conn:   conn,
	}, nil
}

func (pc *PodmanController) exist(imageName string) (bool, error) {
	return containers.Exists(pc.conn, imageName, nil)
}

func (pc *PodmanController) run(imageName string, name string) error {

	// Checks if the container is currently running
	e, _ := pc.status(name) // It is fine for status to error, since it is a health check so erroring means that it fails to see the container running.

	// Abort when there is already a container running, this is a user issue if they try to run the container again when it is already up
	if e {
		fmt.Println("Noticed that container is already running, aborting...")
	}

	// Checks if there is already a container that is in the local storage, we can just run that when needed
	e, err := pc.exist(name)
	if err != nil {
		return err
	}

	if e {
		fmt.Printf("Noticed that container already exist, re-starting the container %s\n", name)
		err = containers.Start(pc.conn, name, nil)
		if err != nil {
			return err
		}
		return nil
	}

	s := specgen.NewSpecGenerator(imageName, false)
	// Config SpecGen
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
	if err = containers.Start(pc.conn, respond.ID, nil); err != nil {
		return err
	}
	fmt.Printf("Container started with\nID: %s\nName:%s\n", respond.ID, name)
	fmt.Println("Access excalidraw on http://localhost:5000/")
	return nil
}

func (pc *PodmanController) stop(name string) error {

	s, err := pc.status(name)

	if !s {
		fmt.Println("Noticed that container is already stopped, aborting...")
		return nil
	}

	stopOtp := &containers.StopOptions{}
	// What is this seconds ? minutes ? hours ?
	// Podman docs kind of sucks
	stopOtp = stopOtp.WithTimeout(1)
	err = containers.Stop(pc.conn, name, stopOtp)
	if err != nil {
		return err
	}
	fmt.Printf("Container stopped with Name: %s\n", name)
	return nil
}

func (pc *PodmanController) status(name string) (bool, error) {
	_, err := containers.RunHealthCheck(pc.conn, name, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (pc *PodmanController) update(imageName string, name string) error {
	status, _ := pc.status(name)

	if status {
		fmt.Println("Noticed that container is already running, stop container before updating image aborting...")
	}

	exist, err := pc.exist(name)
	if err != nil {
		return err
	}
	// If the container already exist remove it
	if exist {
		fmt.Println("Found already exists container")
		rOpt := &containers.RemoveOptions{}
		rOpt = rOpt.WithIgnore(true)
		_, err := containers.Remove(pc.conn, name, rOpt)
		if err != nil {
			return err
		}
		fmt.Println("Container removed")
	}

	img, err := images.Pull(pc.conn, imageName, nil)
	if err != nil {
		return err
	}
	fmt.Println("Pulling image")
	fmt.Printf("Pulling image: %s\n", img)
	return nil
}
