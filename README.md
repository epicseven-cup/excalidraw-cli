# excalidraw-cli
a excalidraw cli that gives you easeier and better controler over a locally excalidraw instance

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
