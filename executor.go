package magicland

type offer struct {
	hostname      string
	port          int
	cpu, memoryMB float64
}

func canHoldOffer(cpu, memoryMB float64) bool {
	return true
}

func localhostOffer(cpu, memoryMB float64) offer {
	return offer{
		hostname: "localhost",
		port:     9000,
		cpu:      cpu,
		memoryMB: memoryMB,
	}
}

func getOffers(cpu, memoryMB float64) ([]offer, error) {
	return []offer{localhostOffer(cpu, memoryMB)}, nil
}
