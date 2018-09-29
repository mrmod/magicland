package magicland

import "testing"

func TestIsRunning(t *testing.T) {
	// Explosion test
	yes, containerID := isRunning("testService")
	_ = yes
	_ = containerID
}
