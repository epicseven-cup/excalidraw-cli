package main

import "flag"

// Flags
var start bool
var status bool
var update bool
var exit bool
var config string

func init() {
	flag.BoolVar(&start, "st", false, "start excalidraw client (shorthand)")
	flag.BoolVar(&start, "start", false, "start excalidraw client")
	flag.BoolVar(&status, "su", false, "excalidraw client status (shorthand)")
	flag.BoolVar(&status, "status", false, "excalidraw client status")
	flag.BoolVar(&update, "u", false, "update excalidraw client image (shorthand)")
	flag.BoolVar(&update, "update", false, "update excalidraw client image")
	flag.BoolVar(&exit, "e", false, "excalidraw client exit")
	flag.BoolVar(&exit, "exit", false, "excalidraw client exit (shorthand)")
	flag.StringVar(&config, "c", "~/config/exclidraw-cli/config", "excalidraw-cli config file")
	flag.StringVar(&config, "config", "~/config/exclidraw-cli/config", "excalidraw-cli config file (shorthand)")
	flag.Parse()
}

func main() {

}
