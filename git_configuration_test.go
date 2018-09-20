package magicland

import "testing"

func TestNonCombustible(t *testing.T) {
	c := NewGitConfiguration("user", "branch", "repoUrl", "serviceName")

	if c.BranchName != "branch" {
		t.Fatal("Expected a Git Configuration")
	}
}
