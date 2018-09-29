package magicland

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
)

// Find a place to run things
func getOffer(cpu, memoryMB float64) (offer, bool) {
	offers, err := getOffers(cpu, memoryMB)
	if err != nil {
		log.Printf("No offers for %2f of CPU and %2fMB RAM", cpu, memoryMB)
		return offer{}, false
	}
	if len(offers) < 1 {
		return offer{}, false
	}
	return offers[0], true
}

// ScheduleService Locates the named service, builds a runtime,
// build a container, then funs the container
func ScheduleService(serviceName string) error {
	gitConfiguration, err := selectServiceByName(serviceName)
	if err != nil {
		log.Printf("%s service does not exist", serviceName)
		return err
	}
	schedulingOffer, ok := getOffer(0.1, 128.0)
	if !ok {
		return fmt.Errorf("Unable to find offers for %s", serviceName)
	}
	app, err := buildAppRuntime(
		schedulingOffer.ipAddress,
		schedulingOffer.port,
		gitConfiguration,
	)
	if err != nil {
		return fmt.Errorf("Unable to build app runtime %s", err.Error())
	}
	ctx := context.Background()
	// Write the express wrapper to the service stage root
	if err := stageContainerApp(app); err != nil {
		return fmt.Errorf("Unable to stage expressified application %s", err.Error())
	}
	// Build a docker container of the app
	runnableContainer, err := buildContainer(ctx, app.RuntimeConfiguration)
	if err != nil {
		return fmt.Errorf("Unable to build a container %s", err.Error())
	}
	runningContainer, err := runContainer(ctx, *runnableContainer)
	if err != nil {
		return fmt.Errorf("Unable to run container %s", err.Error())
	}
	log.Printf("Service %s started in container %s",
		runningContainer.ServiceName,
		runningContainer.ID,
	)
	return nil
}
