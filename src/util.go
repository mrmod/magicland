package magicland

import (
	"fmt"
	"os"
	"strconv"
)

// envOr Returns the environment variable value or a default
func envOr(name, def string) string {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	return v
}

// envOrI Returns the environment variable or default as an Int
func envOrI(name string, def int) int {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(fmt.Sprintf("Invalid configuration for %s: %v", name, err))
	}
	return i
}
