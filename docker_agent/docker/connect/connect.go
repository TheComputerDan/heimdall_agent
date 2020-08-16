package connect

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Containers instantiates a connection to the containers endpoint and returns a list of type `[]types.Container`
func Containers() []types.Container {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}
	return containers
}

// Images instantiates a connection to the images endpoint and returns a list of type `[]types.ImageSummary`
func Images() []types.ImageSummary {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		panic(err)
	}
	return images
}
