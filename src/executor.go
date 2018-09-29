package magicland

import (
	"runtime"

	"github.com/shirou/gopsutil/mem"
)

type offer struct {
	ipAddress, hostname string
	port                int
	cpu, memoryMB       float64
}

func canHoldOffer(cpu, memoryMB float64) bool {
	vmm, _ := mem.VirtualMemory()
	freeMB := float64(vmm.Available) / 1024 / 1024
	if freeMB < memoryMB {
		return false
	}
	if float64(runtime.NumCPU()) < cpu {
		return false
	}
	return true
}

func localhostOffer(cpu, memoryMB float64) offer {
	return offer{
		ipAddress: "127.0.0.1",
		hostname:  "localhost",
		port:      9000,
		cpu:       cpu,
		memoryMB:  memoryMB,
	}
}

func getOffers(cpu, memoryMB float64) ([]offer, error) {
	return []offer{localhostOffer(cpu, memoryMB)}, nil
}
