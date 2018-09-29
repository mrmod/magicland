package magicland

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func init() {
	serviceStageRoot := "/tmp/testServiceName/stage"
	serviceName := "testServiceName"
	_ = os.RemoveAll("/tmp/testServiceName")

	if err := os.MkdirAll(serviceStageRoot, 0755); err != nil {
		fmt.Println("Failed to create", serviceStageRoot, ":", err)
	}
	if yes, containerID := isRunning(serviceName); yes {
		if err := terminate(containerID); err != nil {
			panic(fmt.Sprintf("Failed to terminate %s, %v\n", containerID, err))
		}
	}
}

func TestStageContainer(t *testing.T) {
	rtConfig := RuntimeConfiguration{ServiceName: "testServiceName"}
	dockerClient := new(client.Client)
	builtContainer := container.ContainerCreateCreatedBody{}

	result := rtConfig.stageContainer(builtContainer, dockerClient)
	if result.RuntimeConfiguration.ServiceName != rtConfig.ServiceName {
		t.Fatalf("Expected an encapsulated runtime configuration for %s\n", rtConfig.ServiceName)
	}
}

func TestRunContainer(t *testing.T) {
	serviceStageRoot := "/tmp/testService/stage"
	ctx := context.Background()
	rtConfig := RuntimeConfiguration{
		ServiceName:      "testServiceName",
		ServiceStageRoot: serviceStageRoot,
		entryCommand:     []string{"echo", "container up"},
	}

	runnableContainer, err := buildContainer(ctx, rtConfig)
	if err != nil {
		t.Fatal("Expected a staged container, got", err)
	}

	runningContainer, err := runContainer(ctx, *runnableContainer)
	if err != nil {
		t.Fatal("Expected a running container, got", err)
	}
	fmt.Println("Started ", runningContainer.StagedContainer.ID)
	time.Sleep(3)
	rc, err := runningContainer.DockerClient.ContainerLogs(
		ctx,
		runningContainer.StagedContainer.ID,
		types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
		},
	)
	if err != nil {
		t.Fatal("Expected logs to be readable")
	}
	// Copy the container output to a buffer/ io.Writer
	buf := new(bytes.Buffer)
	io.Copy(buf, rc)

	b, _ := ioutil.ReadAll(buf)
	if matched, _ := regexp.Match("container up", b); !matched {
		t.Fatal("Expected to read 'container up', got", string(b))
	}

}
