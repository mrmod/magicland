package magicland

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

const (
	defaultImage     = "docker.io/library/node:8"
	defaultImageName = "node:8"

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

// Write the Express wrapper app to the stage path as
// magicland.serviceName.js
func stageContainerApp(app appRuntime) error {
	f, err := os.Stat(app.RuntimeConfiguration.ServiceStageRoot)
	if err != nil || !f.IsDir() {
		return fmt.Errorf("%s is not a directory", app.RuntimeConfiguration.ServiceStageRoot)
	}
	expressAppFile := fmt.Sprintf("%s/magicland.%s.js",
		app.RuntimeConfiguration.ServiceStageRoot,
		app.RuntimeConfiguration.ServiceName,
	)

	return ioutil.WriteFile(
		expressAppFile,
		[]byte(app.ExpressApp),
		0755,
	)
}

func buildContainer(ctx context.Context, rtConfig RuntimeConfiguration) (*RunnableContainer, error) {
	var err error
	// Default the container entrypoint
	if len(rtConfig.entryCommand) == 0 {
		rtConfig.entryCommand = defaultEntryCommand
	}
	runnableContainer := &RunnableContainer{
		RuntimeConfiguration: rtConfig,
	}

	runnableContainer.DockerClient, err = client.NewEnvClient()
	if err != nil {
		return runnableContainer, err
	}
	imageReader, err := runnableContainer.DockerClient.ImagePull(
		ctx,
		defaultImage,
		types.ImagePullOptions{},
	)
	if err != nil {
		return runnableContainer, err
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
		return runnableContainer, err
	}

	return runnableContainer, nil
}

// Remove Remove a container
func (this RunnableContainer) Remove(ctx context.Context) error {
	return this.DockerClient.ContainerRemove(
		ctx,
		this.StagedContainer.ID,
		types.ContainerRemoveOptions{},
	)
}
