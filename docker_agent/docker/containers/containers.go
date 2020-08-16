package containers

import (
	"fmt"
	"github.com/docker/docker/api/types"
)

// Containers Returns JSON with the
func List(containers []types.Container) {
	for _, container := range containers {
		fmt.Printf("%s %s %s\n", container.ID[:10], container.Image, container.Names[0])
		fmt.Println(container.Ports)
		fmt.Println(container.Status)
		fmt.Println(container.State)
	}
}
