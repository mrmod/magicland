package magicland

import "testing"

func TestIsRunning(t *testing.T) {
	// Explosion test
	yes, containerID := isRunning("testServiceName")
	_ = yes
	_ = containerID
}
