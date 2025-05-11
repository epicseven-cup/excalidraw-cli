package pkg

import (
	"errors"
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
	exist(name string) (bool, error)
	run(name string) error
	exit(name string) error // Feel like the naming here should be exit since the command is named exit, but stop sounds much better
	status(name string) (bool, error)
	update(name string) error
}

func DetermineEngine() (ContainerEngine, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return DockerEngine, nil
	}

	if _, err := exec.LookPath("podman"); err != nil {
		return PodmanEngine, nil
	}

	return NoEngine, errors.New("no suitable engine found")
}

func NewController(engine ContainerEngine) (*EngineController, error) {
	e, err := DetermineEngine()
	if err != nil {
		return nil, err
	}

	if e == PodmanEngine {
		controller, err := NewPodmanController()

		if err != nil {
			return nil, err
		}

		return &EngineController{
			Engine:     PodmanEngine,
			Controller: controller,
		}, nil
	}

	return nil, nil
}

func (e *EngineController) exist(name string) (bool, error) {
	return e.Controller.exist(name)
}

func (e *EngineController) run(name string) error {
	return e.Controller.run(name)
}

func (e *EngineController) exit(name string) error {
	return e.Controller.exit(name)
}

func (e *EngineController) update(name string) error {
	return e.Controller.update(name)
}

func (e *EngineController) status(name string) (bool, error) {
	return e.Controller.status(name)
}
