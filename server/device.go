package server

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/mush1e/netmon-stack/proto"
)

type Device struct {
	name       string
	interfaces map[string]*proto.InterfaceCounters
	mutex      sync.RWMutex
}

func NewDevice(deviceName string) *Device {
	return &Device{
		name:       deviceName,
		interfaces: make(map[string]*proto.InterfaceCounters),
	}
}

func (d *Device) AddInterface(interfaceName string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if _, ok := d.interfaces[interfaceName]; ok {
		return errors.New("interface already exists")
	}
	d.interfaces[interfaceName] = &proto.InterfaceCounters{
		InterfaceName: interfaceName,
		BytesRx:       0,
		BytesTx:       0,
		PacketsRx:     0,
		PacketsTx:     0,
		Timestamp:     int32(time.Now().UnixMilli()),
	}
	return nil
}

// simulating random network traffic on
// all interfaces for a device
func (d *Device) UpdateCounters() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	for _, iface := range d.interfaces {
		bytesRx := rand.Intn(351) + 50
		bytesTx := rand.Intn(351) + 50
		packetsRx := rand.Intn(21) + 5
		packetsTx := rand.Intn(21) + 5

		iface.BytesRx += int64(bytesRx)
		iface.BytesTx += int64(bytesTx)
		iface.PacketsRx += int64(packetsRx)
		iface.PacketsTx += int64(packetsTx)

		iface.Timestamp = int32(time.Now().UnixMilli())
	}
	return nil
}

func (d *Device) GetCounters(interfaceName string) (*proto.InterfaceCounters, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if iface, ok := d.interfaces[interfaceName]; ok {
		// Return a copy, since we don't want anyone to modify this!
		return &proto.InterfaceCounters{
			InterfaceName: iface.InterfaceName,
			BytesRx:       iface.BytesRx,
			BytesTx:       iface.BytesTx,
			PacketsRx:     iface.PacketsRx,
			PacketsTx:     iface.PacketsTx,
			Timestamp:     iface.Timestamp,
		}, nil
	}
	return nil, errors.New("invalid interface name")
}
