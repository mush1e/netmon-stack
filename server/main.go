package main

import (
	"fmt"
	"sync"
	"time"
)

type server struct {
	devices map[string]*Device
	mutex   sync.RWMutex
}

func main() {
	// Create a device
	device := NewDevice("router1")

	device.AddInterface("eth0")
	device.AddInterface("eth1")

	// Will work on gRPC server stuff later
	// server := &server{
	// 	devices: map[string]*Device{
	// 		"router1": device,
	// 		"switch1": NewDevice("switch1"),
	// 	},
	// }

	device.StartCounterUpdates()

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)

		counters, _ := device.GetCounters("eth0")
		fmt.Printf("eth0 - RX: %d bytes, TX: %d bytes\n",
			counters.BytesRx, counters.BytesTx)
	}
}
