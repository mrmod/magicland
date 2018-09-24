package magicland

import (
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

const (
	defaultImage     = "docker.io/library/alpine"
	defaultImageName = "alpine"

	appRoot = "/app"
)

var defaultEntryCommand = []string{"magicland.init"}

// RunnableContainer is a built container awaiting execution
// on a context
type RunnableContainer struct {
	StagedContainer container.ContainerCreateCreatedBody
	RuntimeConfiguration
	DockerClient *client.Client
}

func buildContainer(rtConfig RuntimeConfiguration) (context.Context, *RunnableContainer, error) {
	var err error
	// Default the container entrypoint
	if len(rtConfig.entryCommand) == 0 {
		rtConfig.entryCommand = defaultEntryCommand
	}
	runnableContainer := &RunnableContainer{
		RuntimeConfiguration: rtConfig,
	}
	ctx := context.Background()
	runnableContainer.DockerClient, err = client.NewEnvClient()
	if err != nil {
		return ctx, runnableContainer, err
	}
	imageReader, err := runnableContainer.DockerClient.ImagePull(
		ctx,
		defaultImage,
		types.ImagePullOptions{},
	)
	if err != nil {
		return ctx, runnableContainer, err
	}
	// Relay status output to Stdout
	io.Copy(os.Stdout, imageReader)
	containerConfig := &container.Config{
		Image: defaultImageName,
		Cmd:   runnableContainer.RuntimeConfiguration.entryCommand,
		Tty:   true,
	}
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			mount.Mount{
				Type:   mount.TypeBind,
				Source: runnableContainer.RuntimeConfiguration.ServiceStageRoot,
				Target: appRoot,
			},
		},
	}
	// Create a container with no host or network configuration
	runnableContainer.StagedContainer, err = runnableContainer.DockerClient.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig, // Host config
		nil,        // no Network config
		runnableContainer.RuntimeConfiguration.ServiceName)

	if err != nil {
		return ctx, runnableContainer, err
	}

	return ctx, runnableContainer, nil
}

func (this RunnableContainer) Remove(ctx context.Context) error {
	return this.DockerClient.ContainerRemove(
		ctx,
		this.StagedContainer.ID,
		types.ContainerRemoveOptions{},
	)
}
