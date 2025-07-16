package main

import (
	"fmt"
	"time"
)

func main() {
	// Create a device
	device := NewDevice("router1")

	device.AddInterface("eth0")
	device.AddInterface("eth1")

	device.StartCounterUpdates()

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)

		counters, _ := device.GetCounters("eth0")
		fmt.Printf("eth0 - RX: %d bytes, TX: %d bytes\n",
			counters.BytesRx, counters.BytesTx)
	}
}
