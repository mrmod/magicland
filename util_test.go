package magicland

import (
	"os"
	"testing"
)

func TestEnvOrI(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("OK: PATH cannot be coerced to Int")
		}
	}()
	if i := envOrI("PATH", 42); i != 42 {
		t.Fatalf("Expected 42, got %d", i)
	}
	// I wonder which test systems are more than a shell deep where
	// this might fail. Just in case...
	// actualShellLevel, _ := strconv.Atoi(os.Getenv("SHLVL"))
	// Common to most Unix shells
	if i := envOrI("SHLVL", 42); i != 1 {
		t.Fatalf("Expected %s, got %d", os.Getenv("SHLVL"), i)
	}
}
