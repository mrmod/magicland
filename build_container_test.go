package magicland

import (
	"fmt"
	"os"
	"regexp"
	"testing"
)

var serviceStageRoot = "/tmp/testService/stage"

func init() {

	_ = os.RemoveAll("/tmp/testService")

	if err := os.MkdirAll(serviceStageRoot, 0755); err != nil {
		fmt.Println("Failed to create", serviceStageRoot, ":", err)
	}
}

func TestBuildContainer(t *testing.T) {
	rtConfig := RuntimeConfiguration{
		ServiceName:      "testService",
		ServiceStageRoot: serviceStageRoot,
	}
	if yes, containerID := isRunning(rtConfig.ServiceName); yes {
		if err := terminate(containerID); err != nil {
			t.Fatalf("Failed to terminate %s, %v\n", containerID, err)
		}
	}
	ctx, runnableContainer, err := buildContainer(rtConfig)
	defer func() {
		_ = runnableContainer.Remove(ctx)
	}()
	if err != nil {
		if m, _ := regexp.Match("Is the docker daemon running", []byte(err.Error())); m {
			t.Skipf("Expected the docker daemon to be running\n")
		}
		t.Fatal("Expected a new container created body, got", err)
	}
	_ = ctx
	if len(runnableContainer.StagedContainer.Warnings) != 0 {
		t.Fatal("Expected no warnings creating the container, got", runnableContainer.StagedContainer.Warnings)
	}
	if err := ctx.Err(); err != nil {
		t.Fatal("Expected no errors on context, got", err)
	}
}
