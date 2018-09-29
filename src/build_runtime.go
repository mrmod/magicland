package magicland

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
)

var (
	defaultCORSHandler = `
const cors = (response) => {
  response.set("Access-Control-Allow-Origin", "*");
  response.set("Access-Control-Allow-Methods", "*");
  response.set("Access-Control-Allow-Headers", "Content-Type");
  response.set("Access-Control-Max-Age", "3600");
  return response
}
`

	defaultServiceStartedNotifier = `
const notifyServiceStarted = (rtPort, rtHost, serviceName) => {
	console.log("Started", serviceName, "on", rtHost +":"+rtPort)
}
`
)

// RuntimeConfiguration Defines where a service will run
type RuntimeConfiguration struct {
	// TODO : Rename to IPAddress
	Host                   string
	Port                   int
	CORSHandler            string // Function definition in JS
	ServiceName            string // Magicland service name
	ServiceStageRoot       string
	ServiceStartedNotifier string   // Function definition in JS
	entryCommand           []string // Command to run on container execution
}

// appRuntime contains the git and runtime configuration
type appRuntime struct {
	RuntimeConfiguration
	GitConfiguration
	ExpressApp string
}

func newRuntimeConfiguration(host string, port int, gitConfig GitConfiguration) RuntimeConfiguration {
	serviceStageRoot := gitConfig.servicePath()
	return RuntimeConfiguration{
		Host:                   host,
		Port:                   port,
		CORSHandler:            defaultCORSHandler,
		ServiceName:            gitConfig.ServiceName,
		ServiceStageRoot:       serviceStageRoot,
		ServiceStartedNotifier: defaultServiceStartedNotifier,
	}
}

func buildAppRuntime(host string, port int, gitConfig GitConfiguration) (appRuntime, error) {
	var serviceConfiguration string
	rtConfig := newRuntimeConfiguration(host, port, gitConfig)
	serviceConfiguration = buildExpressConfiguration(rtConfig)
	if serviceConfiguration == "" {
		return appRuntime{}, errors.New("Unable to build service configuration")
	}
	return appRuntime{
		rtConfig,
		gitConfig,
		serviceConfiguration,
	}, nil
}

func buildExpressConfiguration(rtConfig RuntimeConfiguration) string {
	// TODO: Cors for public DNS configuration
	rawTemplate := `
const app = require('express')()
const {handle} = require('./index')
// CORS Handler
{{.CORSHandler}}
// ServiceStartedNotifier
{{.ServiceStartedNotifier}}
const port = {{.Port}}
const host = '{{.Host}}'
const serviceName = '{{.ServiceName}}'
const terminator = (ec) => {
	// 0 is the same as null/undefined
	if (ec) {
		process.exit(ec)
	}
	process.exit(0)
}
app.all("/", (req, res) => {
	try {
		handle(req, res, terminator)
	} catch(e) {
		console.log(e)
		terminator(1)
	}
	terminator(0)
})
app.listen(port, notifyServiceStarted(port, host, serviceName))
`

	t, err := template.New("serviceHandler").Parse(rawTemplate)
	if err != nil {
		fmt.Println("Error compiling template", err)
		return ""
	}
	buf := &bytes.Buffer{}
	t.Execute(buf, rtConfig)
	return buf.String()
}
