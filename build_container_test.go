package magicland

import (
	"regexp"
	"testing"
)

func TestBuildContainer(t *testing.T) {
	rtConfig := RuntimeConfiguration{ServiceName: "testService"}
	ctx, container, err := buildContainer(rtConfig)
	if err != nil {
		if m, _ := regexp.Match("Is the docker daemon running", []byte(err.Error())); m {
			t.Skipf("Expected the docker daemon to be running\n")
		}
		t.Fatal("Expected a new container created body, got", err)
	}
	_ = ctx
	if len(container.Warnings) != 0 {
		t.Fatal("Expected no warnings creating the container, got", container.Warnings)
	}
	if err := ctx.Err(); err != nil {
		t.Fatal("Expected no errors on context, got", err)
	}
}
