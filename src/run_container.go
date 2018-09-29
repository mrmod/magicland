package magicland

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// RunningContainer is the result of a RunnableContainer
// execution
type RunningContainer struct {
	ID string
	RunnableContainer
	Stdout io.ReadWriter
	Stderr io.ReadWriter
	Stdin  io.Writer
}

func (this RuntimeConfiguration) stageContainer(builtContainer container.ContainerCreateCreatedBody, dockerClient *client.Client) *RunnableContainer {
	return &RunnableContainer{
		StagedContainer:      builtContainer,
		RuntimeConfiguration: this,
		DockerClient:         dockerClient,
	}
}

func runContainer(ctx context.Context, rc RunnableContainer) (RunningContainer, error) {
	runningContainer := RunningContainer{
		RunnableContainer: rc,
	}
	err := rc.DockerClient.ContainerStart(
		ctx,
		rc.StagedContainer.ID,
		types.ContainerStartOptions{},
	)
	if err != nil {
		return runningContainer, err
	}

	return runningContainer, nil
}
