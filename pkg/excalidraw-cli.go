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
	// Checks for image exist
	exist(imageName string) (bool, error)
	// Runs the container
	run(imageName string, name string) error
	// Exits the container
	exit(name string) error // Feel like the naming here should be exit since the command is named exit, but stop sounds much better
	// Check the status of the container
	status(name string) (bool, error)
	// Update the image
	update(imageName string) error
}

func DetermineEngine() (ContainerEngine, error) {
	if _, err := exec.LookPath("docker"); err == nil {
		fmt.Println("Determined Docker Engine")
		return DockerEngine, nil
	}

	if _, err := exec.LookPath("podman"); err == nil {
		fmt.Println("Determined Podman Engine")
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

	//if e == DockerEngine {
	//	controller, err := NewDockerController(system)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return &EngineController{
	//		Engine:     DockerEngine,
	//		Controller: controller,
	//	}, nil
	//}

	return nil, nil
}

func (e *EngineController) Exist(name string) (bool, error) {
	return e.Controller.exist(name)
}

func (e *EngineController) Run(imageName string, name string) error {
	return e.Controller.run(imageName, name)
}

func (e *EngineController) Exit(name string) error {
	return e.Controller.exit(name)
}

func (e *EngineController) Update(name string) error {
	return e.Controller.update(name)
}

func (e *EngineController) Status(name string) (bool, error) {
	return e.Controller.status(name)
}
