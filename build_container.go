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
	defaultImage       = "docker.io/library/alpine"
	defaultImageName   = "alpine"
	defaultInitCommand = "magicland.init"
	appRoot            = "/app"
)

func buildContainer(rtConfig RuntimeConfiguration) (context.Context, container.ContainerCreateCreatedBody, error) {
	var newContainer container.ContainerCreateCreatedBody
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {

		return ctx, newContainer, err
	}
	imageReader, err := cli.ImagePull(ctx, defaultImage, types.ImagePullOptions{})
	if err != nil {
		return ctx, newContainer, err
	}
	// Relay status output to Stdout
	io.Copy(os.Stdout, imageReader)
	containerConfig := &container.Config{
		Image: defaultImageName,
		Cmd:   []string{defaultInitCommand},
		Tty:   true,
	}
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			mount.Mount{
				Type:   mount.TypeBind,
				Source: rtConfig.ServiceStageRoot,
				Target: appRoot,
			},
		},
	}
	// Create a container with no host or network configuration
	newContainer, err = cli.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig, // Host config
		nil,        // no Network config
		rtConfig.ServiceName)

	if err != nil {
		return ctx, newContainer, err
	}

	return ctx, newContainer, nil
}
