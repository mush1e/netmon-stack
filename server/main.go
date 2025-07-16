package main

import (
	"errors"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/mush1e/netmon-stack/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	// For default stuff belonging to NetworkTelemetry server
	proto.UnimplementedNetworkTelemetryServer

	// devices and mtx
	devices map[string]*Device
	mutex   sync.RWMutex
}

func (s *server) Subscribe(req *proto.SubscribeRequest, stream proto.NetworkTelemetry_SubscribeServer) error {
	nameSlice := strings.Split(req.InterfaceName, ":")
	if len(nameSlice) != 2 {
		return stream.Send(&proto.SubscribeResponse{
			Response: &proto.SubscribeResponse_Error{
				Error: &proto.Error{
					Code:    proto.ErrorCode_DOES_NOT_EXIST,
					Message: "invalid interface name format; expected DEVICE:INTERFACE",
				},
			},
			ResponseTimestamp: time.Now().UnixMilli(),
		})
	}

	deviceName := nameSlice[0]
	ifaceName := nameSlice[1]

	device := s.getDevice(deviceName)
	if device == nil {
		return stream.Send(&proto.SubscribeResponse{
			Response: &proto.SubscribeResponse_Error{
				Error: &proto.Error{
					Code:    proto.ErrorCode_DOES_NOT_EXIST,
					Message: "device '" + deviceName + "' not found",
				},
			},
			ResponseTimestamp: time.Now().UnixMilli(),
		})
	}

	interval := time.Duration(req.IntervalMs) * time.Millisecond
	if interval <= 0 {
		interval = 1000 * time.Millisecond // default to 1s
	}

	switch req.Mode {
	case proto.SubscriptionMode_ONCE:
		ifaceCounters, err := device.GetCounters(ifaceName)
		if err != nil {
			return stream.Send(&proto.SubscribeResponse{
				Response: &proto.SubscribeResponse_Error{
					Error: &proto.Error{
						Code:    proto.ErrorCode_NOT_ACTIVE,
						Message: "failed to get counters for interface '" + ifaceName + "': " + err.Error(),
					},
				},
				ResponseTimestamp: time.Now().UnixMilli(),
			})
		}

		return stream.Send(&proto.SubscribeResponse{
			Response: &proto.SubscribeResponse_Counters{
				Counters: ifaceCounters,
			},
			ResponseTimestamp: time.Now().UnixMilli(),
		})

	case proto.SubscriptionMode_STREAM:
		for {
			select {
			case <-stream.Context().Done():
				return nil
			default:
				ifaceCounters, err := device.GetCounters(ifaceName)
				if err != nil {
					_ = stream.Send(&proto.SubscribeResponse{
						Response: &proto.SubscribeResponse_Error{
							Error: &proto.Error{
								Code:    proto.ErrorCode_NOT_ACTIVE,
								Message: "stream error getting counters for '" + ifaceName + "': " + err.Error(),
							},
						},
						ResponseTimestamp: time.Now().UnixMilli(),
					})
					return nil // End stream on error
				}

				if err := stream.Send(&proto.SubscribeResponse{
					Response: &proto.SubscribeResponse_Counters{
						Counters: ifaceCounters,
					},
					ResponseTimestamp: time.Now().UnixMilli(),
				}); err != nil {
					return err
				}

				time.Sleep(interval)
			}
		}

	case proto.SubscriptionMode_POLL:
		return stream.Send(&proto.SubscribeResponse{
			Response: &proto.SubscribeResponse_Error{
				Error: &proto.Error{
					Code:    proto.ErrorCode_NOT_ACTIVE,
					Message: "POLL mode is not implemented yet",
				},
			},
			ResponseTimestamp: time.Now().UnixMilli(),
		})

	default:
		return stream.Send(&proto.SubscribeResponse{
			Response: &proto.SubscribeResponse_Error{
				Error: &proto.Error{
					Code:    proto.ErrorCode_DOES_NOT_EXIST,
					Message: "invalid subscription mode",
				},
			},
			ResponseTimestamp: time.Now().UnixMilli(),
		})
	}
}

func NewServer() *server {
	return &server{
		devices: make(map[string]*Device),
	}
}

func (s *server) getDevice(deviceName string) *Device {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if _, ok := s.devices[deviceName]; ok {
		return s.devices[deviceName]
	}
	return nil
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

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcServer := grpc.NewServer()
	proto.RegisterNetworkTelemetryServer(grpcServer, srv)
	reflection.Register(grpcServer)

	log.Println("Starting gRPC server on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err.Error())
	}

}
