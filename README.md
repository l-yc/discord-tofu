# discord-tofu

A discord bot written in Go that uses [tofu-ai](https://github.com/wernjie/tofu-ai).

## Makefile Tasks
- `make docker`: builds the docker container
- `make run`: runs the built docker container
- `make build`: builds the Go binary
- `make pack`: creates a zip archive which for deploying directly
- `make install`: installs the needed python libraries on the deployment machine (requires python3 to be installed)
