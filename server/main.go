package main

import (
	"errors"
	"log"
	"net"
	"sync"

	"github.com/mush1e/netmon-stack/proto"
	"google.golang.org/grpc"
)

type server struct {
	// For default stuff belonging to NetworkTelemetry server
	proto.UnimplementedNetworkTelemetryServer

	// devices and mtx
	devices map[string]*Device
	mutex   sync.RWMutex
}

func (s *server) Subscribe(req *proto.SubscribeRequest, stream proto.NetworkTelemetry_SubscribeServer) error {
	// streaming logic
	return nil
}

func NewServer() *server {
	return &server{
		devices: make(map[string]*Device),
	}
}

func (s *server) addDevice(deviceName string) (*Device, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.devices[deviceName]; ok {
		return nil, errors.New("device already exists on server")
	}
	s.devices[deviceName] = NewDevice(deviceName)
	return s.devices[deviceName], nil
}

func main() {
	srv := NewServer()
	if dev, err := srv.addDevice("router1"); err == nil {
		dev.AddInterface("eth0")
		dev.AddInterface("eth1")
		dev.AddInterface("eth2")
	} else {
		log.Printf("router1 already exists")
	}

	if dev, err := srv.addDevice("switch1"); err == nil {
		dev.AddInterface("gi0/0")
		dev.AddInterface("gi0/1")
	} else {
		log.Printf("switch1 already exists")
	}

	srv.mutex.RLock()
	for _, device := range srv.devices {
		device.StartCounterUpdates()
	}
	srv.mutex.RUnlock()

	listener, err := net.Listen("tcp", "500051")
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcServer := grpc.NewServer()
	proto.RegisterNetworkTelemetryServer(grpcServer, srv)

	log.Println("Starting gRPC server on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err.Error())
	}

}
