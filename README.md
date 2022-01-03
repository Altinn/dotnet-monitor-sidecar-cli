# dotnet-monitor-sidecar-cli
[![Go Report Card](https://goreportcard.com/badge/github.com/altinn/dotnet-monitor-sidecar-cli)](https://goreportcard.com/report/github.com/altinn/dotnet-monitor-sidecar-cli)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/altinn/dotnet-monitor-sidecar-cli.svg)](https://github.com/altinn/dotnet-monitor-sidecar-cli) 

CLI for adding [dotnet-monitor](https://github.com/dotnet/dotnet-monitor) sidecar to pod in Kubernetes.

The main goal of this cli is to ease the process of adding the dotnet-monitor docker image with existing pods in kubernetes

## Install

Binary downloads of the dmsctl can be found on the [releases page](https://github.com/Altinn/dotnet-monitor-sidecar-cli/releases/latest)

Unpack the dmsctl binary and add it to your PATH and you should be good to go.

## Usage

The goale is to let the cli tool to be self documenting with the `-h` flag.

Examples:
```
dmsctl -h
dmsctl add -h
```

The docs are also available as markdown [here](docs/dmsctl.md)
## Contributing
Please read CONTRIBUTING.md for details on our code of conduct, and the process for submitting pull requests to us.