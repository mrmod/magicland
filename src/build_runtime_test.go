package magicland

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"testing"
	"time"
)

func TestBuildExpressConfiguration(t *testing.T) {
	serviceName := "testServiceName"

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

func TestRuntimeViability(t *testing.T) {
	// Create a fake app to run in the common UTS
	serviceName := "testServiceName"
	app := "node/runtime-viability-test.js"
	gitConfig := GitConfiguration{ServiceName: serviceName}
	rtConfig := newRuntimeConfiguration("host01.local", 8000, gitConfig)
	s := buildExpressConfiguration(rtConfig)

	// Stage the configuration to the filesystem in our ./node scratch space
	if err := ioutil.WriteFile(app, []byte(s), 0755); err != nil {
		t.Fatal("Expected to write a node app configuration, got", err)
	}
	handlerFile := "node/index.js"
	handlerFunction := []byte("exports.handle = (req, res) => res.send('HELLO WORLD')")
	if err := ioutil.WriteFile(handlerFile, handlerFunction, 0755); err != nil {
		t.Fatal("Failed to write test handler", err)
	}

	// Simulate container entry
	c := exec.Command("node", app)
	// Clean-up our process
	defer func() {
		t.Log("Cleaning up process", c.Process.Pid)
		if c.ProcessState.Exited() {
			t.Log("Node application properly exited, no cleanup needed")
			return
		}
		if err := c.Process.Kill(); err != nil {
			t.Fatalf("Failure in shutdown while killing %d, %v\n", c.Process.Pid, err)
		}
		// The "container entrypoint" should exit after serving a single request
		t.Fatal("Node application did not exist after request, cleanup was forced")
	}()
	// Pre-launch HTTP request test
	go func() {
		// Fire a request at the API to trigger its termination
		t.Log("Sending runtimeviability HTTP request")
		// Give time for OSX users to allow the firewall breach
		time.Sleep(5 * time.Second)
		request, err := http.NewRequest("GET", "http://localhost:8000", nil)
		if err != nil {
			t.Fatal("Failure in test HTTP request harness", err)
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Fatal("Failure in test HTTP response harness", err)
		}
		if b, _ := ioutil.ReadAll(response.Body); string(b) != "HELLO WORLD" {
			t.Fatalf("Expected HELLO WORLD, got %s\n", string(b))
		}
	}()
	// Launch the Node handler and block
	if sose, err := c.CombinedOutput(); err != nil {
		t.Fatal("Failed to get combined output in test harness", err, string(sose))
	}

}
