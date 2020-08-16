package images

import (
	"fmt"
	"github.com/docker/docker/api/types"
)

// List returns image information from the host with the agent setup
func List(images []types.ImageSummary) {
	for _, image := range images {
		fmt.Printf("%s", image.RepoTags)
	}
}
