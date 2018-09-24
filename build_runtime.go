package magicland

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
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
	Host                   string
	Port                   int
	CORSHandler            string // Function definition in JS
	ServiceName            string // Magicland service name
	ServiceStageRoot       string
	ServiceStartedNotifier string   // Function definition in JS
	entryCommand           []string // Command to run on container execution
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

func buildRuntime(host string, port int, gitConfig GitConfiguration) error {
	rtConfig := newRuntimeConfiguration(host, port, gitConfig)
	serviceConfiguration := buildExpressConfiguration(rtConfig)
	if serviceConfiguration == "" {
		return errors.New("Unable to build service configuration")
	}
	return nil
}

func buildExpressConfiguration(rtConfig RuntimeConfiguration) string {
	// TODO: Cors for public DNS configuration
	rawTemplate := `
const express = require('express')
const {handle} = require('./index')
// CORS Handler
{{.CORSHandler}}
// ServiceStartedNotifier
{{.ServiceStartedNotifier}}
const app = express()
const port = {{.Port}}
app.all("/", handle)
app.listen(port, notifyServiceStarted({{.Port}}, {{.Host}}, {{.ServiceName}}))
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
