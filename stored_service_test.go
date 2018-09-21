package magicland

import (
	"testing"
)

var testServiceName = "testServiceName"

func TestSelectServiceByName(t *testing.T) {
	_ = saveService(GitConfiguration{"user", "branch", "repo", "testServiceName"})
	config, err := selectServiceByName(testServiceName)
	if err != nil {
		t.Fatal("Expected a GitConfiguration, got ", err)
	}

	if config.ServiceName != testServiceName {
		t.Fatalf("Expected %s, got %s", testServiceName, config.ServiceName)
	}
}

func TestSaveService(t *testing.T) {
	err := saveService(GitConfiguration{"user", "branch", "repo", "testServiceName"})
	if err != nil {
		t.Fatal("Expected no error, got ", err)
	}
}
