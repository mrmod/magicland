package magicland

import "log"

// ExecuteService Locates the named service, builds a runtime,
// build a container, then funs the container
func ExecuteService(serviceName string) error {
	gitConfiguration, err := selectServiceByName(serviceName)
	if err != nil {
		log.Printf("%s service does not exist", serviceName)
		return err
	}
	// Find a place to run things
	cpu := 0.1
	memoryMB := 128.0
	offers, err := getOffers(cpu, memoryMB)
	if err != nil {
		log.Printf("No offers for %2f of CPU and %2fMB RAM", cpu, memoryMB)
		return err
	}
	_ = offers
	_ = gitConfiguration
	return nil
}
