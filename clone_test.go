package magicland

import (
	"os"
	"testing"
)

var publicGitConfig GitConfiguration
var serviceName = "magicland-test-repo"

func init() {
	s, err := os.Stat("/tmp/" + serviceName)
	if err != nil {
		println(err)
	} else {
		if s.IsDir() {
			println("Removing existing clone")
			os.RemoveAll("/tmp/" + serviceName)
		}
	}

	publicGitConfig = NewGitConfiguration(
		"user",
		"master",
		"https://github.com/mrmod/magicland.git",
		serviceName)
}
func TestPublicClone(t *testing.T) {
	t.Log("It should clone a public repository")
	if err := PublicClone(publicGitConfig); err != nil {
		t.Fatalf("Expected to clone %s but failed %v", publicGitConfig.RepositoryURL, err)
	}
}
