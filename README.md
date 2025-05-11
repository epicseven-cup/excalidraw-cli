# excalidraw-cli
a excalidraw cli that gives you easeier and better controler over a locally excalidraw instance


# [Warning] Currently this package only supports Podman, I need to setup some testing for docker to make sure it works fully.

# Install
1. You will need to install go on your machine: https://go.dev/doc/install
2. Setup GOPATH

Add the following to your shell config
```bash
export PATH=${PATH}:$HOME/go/bin
```
More information: https://go.dev/wiki/GOPATH#gopath-variable

3. Install the binary
```bash
go install github.com/epicseven-cup/excalidraw-cli@latest 


4. Make sure you have a supported container engine installed and setup. (e.g. podman, docker...)
```

There could be delays between the Goproxy and GitHub binarys, you can use the direct setup
```bash
GOPROXY=direct go install github.com/epicseven-cup/excalidraw-cli@latest
```

# How to use


```bash
 excalidraw-cli -h
```

# Example

Starts the excalidraw docker container, you can access it on  [localhost:5000](http://localhost:5000)
```bash
excalidraw-cli -start 
```
