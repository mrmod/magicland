package magicland

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func isRunning(serviceName string) (bool, string) {
	var containerID string
	if !strings.HasPrefix(serviceName, "/") {
		serviceName = "/" + serviceName
	}
	ctx := context.Background()
	c, err := client.NewEnvClient()
	if err != nil {
		return false, containerID
	}
	containers, err := c.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return false, containerID
	}
	for _, runningContainer := range containers {
		fmt.Printf("Container %s is in state %s\n", runningContainer.ID, runningContainer.Status)
		for _, name := range runningContainer.Names {
			if name == serviceName {
				return true, runningContainer.ID
			}
		}
	}
	return false, containerID
}

func terminate(containerID string) error {
	ctx := context.Background()
	c, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	return c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{})
}
