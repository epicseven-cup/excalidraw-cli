package main

import (
	"flag"
	"fmt"
	"github.com/epicseven-cup/excalidraw-cli/pkg"
	"runtime"
)

const (
	EXCALIDRAW_IMAGE     = "excalidraw/excalidraw"
	EXCALIDRAW_CONTAINER = "excalidraw-cli-container"
)

// Flags
var start bool
var status bool
var update bool
var stop bool
var config string
var engine string

func init() {
	flag.BoolVar(&start, "st", false, "start excalidraw client (shorthand)")
	flag.BoolVar(&start, "start", false, "start excalidraw client")
	flag.BoolVar(&status, "su", false, "excalidraw client status (shorthand)")
	flag.BoolVar(&status, "status", false, "excalidraw client status")
	flag.BoolVar(&update, "u", false, "update excalidraw client image (shorthand)")
	flag.BoolVar(&update, "update", false, "update excalidraw client image")
	flag.BoolVar(&stop, "sp", false, "excalidraw client exit")
	flag.BoolVar(&stop, "stop", false, "excalidraw client exit (shorthand)")
	flag.StringVar(&config, "c", "~/config/exclidraw-cli/config", "excalidraw-cli config file")
	flag.StringVar(&config, "config", "~/config/exclidraw-cli/config", "excalidraw-cli config file (shorthand)")
	flag.StringVar(&engine, "e", "", "container engine (shorthand)")
	flag.StringVar(&engine, "engine", "", "container engine")
}

func DetermineEngine(system string) (*pkg.EngineController, error) {
	c, err := pkg.NewController(system)

	if err != nil {
		return nil, err
	}

	if engine == "podman" {
		podmanController, err := pkg.NewPodmanController(system)
		if err != nil {
			return nil, err
		}
		fmt.Println("engine flag enable, overriding engine to podman engine")
		c.Controller = podmanController
	}

	if engine == "docker" {
		dockerController, err := pkg.NewDockerController(system)
		if err != nil {
			return nil, err
		}
		fmt.Println("engine flag enable, overriding engine to docker engine")
		c.Controller = dockerController
	}

	return c, nil

}

func excalidraw(system string) error {
	c, err := DetermineEngine(system)
	if err != nil {
		return err
	}

	if start {
		err = c.Run(EXCALIDRAW_IMAGE, EXCALIDRAW_CONTAINER)
		if err != nil {
			return err
		}
	}

	if status {
		b, _ := c.Status(EXCALIDRAW_CONTAINER)
		fmt.Printf("Checking if the container is already running: %t\n", b)
	}

	if update {
		err := c.Update(EXCALIDRAW_IMAGE, EXCALIDRAW_CONTAINER)
		if err != nil {
			return err
		}
	}

	if stop {
		err := c.Stop(EXCALIDRAW_CONTAINER)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	fmt.Println("Running excalidraw-cli")
	flag.Parse()
	fmt.Println("Detecting operating system")
	var goos string
	switch goos = runtime.GOOS; goos {
	case "darwin":
		fmt.Println("OS: darwin / macOS")
	case "linux":
		fmt.Println("OS: linux")
	case "windows":
		fmt.Println("OS: windows")
		fmt.Println("Not supported on Windows")
		return
	default:
		fmt.Println("What is this? Did you really put a monitor and keyboard together and call it a computer?")
		return
	}
	if err := excalidraw(goos); err != nil {
		fmt.Println(err)
	}
}
