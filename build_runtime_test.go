package magicland

import (
	"regexp"
	"testing"
)

func TestBuildExpressConfiguration(t *testing.T) {
	serviceName := "testService"

	gitConfig := GitConfiguration{ServiceName: serviceName}
	rtConfig := newRuntimeConfiguration("host", 8000, gitConfig)
	s := buildExpressConfiguration(rtConfig)

	if len(s) < 10 {
		t.Fatal("Expected a larger configuration than", s)
	}
	if _, err := regexp.MatchString("const cors =", s); err != nil {
		t.Fatal("Expected to find a CORS handler")
	}
	if _, err := regexp.MatchString("const notifyServiceStarted =", s); err != nil {
		t.Fatal("Expected to find a service up notifier")
	}
	if _, err := regexp.MatchString(", notifyServiceStart.?8000", s); err != nil {
		t.Fatal("Expected to find a Express listener callback")
	}
}
