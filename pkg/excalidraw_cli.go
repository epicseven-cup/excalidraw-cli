package pkg

import (
	"errors"
	"fmt"
	"os/exec"
)

type ContainerEngine int

const (
	DockerEngine ContainerEngine = iota
	PodmanEngine
	NoEngine
)

type EngineController struct {
	Engine     ContainerEngine
	Controller Controller
}

type Controller interface {
	// Checks if the container already exist doesn't matter if it is already running or not
	exist(imageName string) (bool, error)
	// Runs the container, checks if there is already a container stopped use that instead
	run(imageName string, name string) error
	// Stops the container
	stop(name string) error // Feel like the naming here should be exited since the command is named exit, but stop sounds much better
	// Check the status of the container is running or not
	status(name string) (bool, error)
	// Update the image, and remove the old image and container
	update(imageName string, name string) error
}

func DetermineEngine() (ContainerEngine, error) {
	if _, err := exec.LookPath("docker"); err == nil {
		fmt.Println("Found Docker Engine")
		return DockerEngine, nil
	}

	if _, err := exec.LookPath("podman"); err == nil {
		fmt.Println("Found Podman Engine")
		return PodmanEngine, nil
	}

	return NoEngine, errors.New("no suitable engine found")
}

func NewController(system string) (*EngineController, error) {
	e, err := DetermineEngine()
	if err != nil {
		return nil, err
	}

	if e == PodmanEngine {
		controller, err := NewPodmanController(system)

		if err != nil {
			return nil, err
		}

		return &EngineController{
			Engine:     PodmanEngine,
			Controller: controller,
		}, nil
	}

	if e == DockerEngine {
		controller, err := NewDockerController(system)
		if err != nil {
			return nil, err
		}
		return &EngineController{
			Engine:     DockerEngine,
			Controller: controller,
		}, nil
	}

	return nil, nil
}

func (e *EngineController) Exist(name string) (bool, error) {
	return e.Controller.exist(name)
}

func (e *EngineController) Run(imageName string, name string) error {
	return e.Controller.run(imageName, name)
}

func (e *EngineController) Stop(name string) error {
	return e.Controller.stop(name)
}

func (e *EngineController) Update(imageName string, name string) error {
	return e.Controller.update(imageName, name)
}

func (e *EngineController) Status(name string) (bool, error) {
	return e.Controller.status(name)
}
