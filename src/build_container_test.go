package magicland

import (
	"fmt"
	"os"
	"regexp"
	"testing"

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

func TestBuildContainer(t *testing.T) {
	serviceStageRoot := "/tmp/testService/stage"
	rtConfig := RuntimeConfiguration{
		ServiceName:      "testService",
		ServiceStageRoot: serviceStageRoot,
	}

	ctx := context.Background()
	runnableContainer, err := buildContainer(ctx, rtConfig)
	defer func() {
		_ = runnableContainer.Remove(ctx)
	}()
	if err != nil {
		if m, _ := regexp.Match("Is the docker daemon running", []byte(err.Error())); m {
			t.Skipf("Expected the docker daemon to be running\n")
		}
		t.Fatal("Expected a new container created body, got", err)
	}

	if len(runnableContainer.StagedContainer.Warnings) != 0 {
		t.Fatal("Expected no warnings creating the container, got", runnableContainer.StagedContainer.Warnings)
	}
	if err := ctx.Err(); err != nil {
		t.Fatal("Expected no errors on context, got", err)
	}
}
