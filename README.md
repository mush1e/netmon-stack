# NetMon Stack - Real-Time Network Monitoring System

A **gRPC-based network telemetry system** that simulates Cisco network devices and provides real-time interface counter streaming. Built to teach modern network monitoring technologies used in production environments.

## Project Overview

**NetMon Stack** is a complete network monitoring solution that demonstrates:
- **gRPC & Protocol Buffers** for high-performance network communication
- **gNMI-style telemetry** for real-time network data streaming  
- **Go concurrency patterns** for handling multiple device simulations
- **Realistic network simulation** with traffic patterns and interface counters
- **Production-ready architecture** with proper error handling and thread safety

## Architecture

```
┌─────────────────┐    gRPC Stream    ┌─────────────────┐
│                 │◄─────────────────►│                 │
│  Client Apps    │   Subscribe()     │  gRPC Server    │
│                 │   device:iface    │                 │
└─────────────────┘                   └─────────────────┘
                                               │
                                               │ manages
                                               ▼
                                      ┌─────────────────┐
                                      │ Device Simulator│
                                      │ ┌─────────────┐ │
                                      │ │ router1     │ │
                                      │ │  └─eth0     │ │
                                      │ │  └─eth1     │ │
                                      │ │  └─eth2     │ │
                                      │ └─────────────┘ │
                                      │ ┌─────────────┐ │
                                      │ │ switch1     │ │
                                      │ │  └─gi0/0    │ │
                                      │ │  └─gi0/1    │ │
                                      │ └─────────────┘ │
                                      └─────────────────┘
```

## Features

### **Phase 1: Protocol Foundation** (COMPLETED)
- [x] **Protocol Buffer definitions** for network telemetry data
- [x] **gNMI-style gRPC service** with Subscribe streaming
- [x] **Device simulation** with realistic interface counters
- [x] **Concurrent counter updates** using goroutines
- [x] **Thread-safe operations** with RWMutex patterns

### **Phase 2: gRPC Server Implementation** (IN PROGRESS)
- [ ] **Server struct** with device management
- [ ] **Subscribe method** implementation  
- [ ] **String parsing** for "device:interface" format
- [ ] **Error handling** for missing devices/interfaces
- [ ] **Multiple subscription modes** (STREAM, ONCE, POLL)

### **Phase 3: Advanced Features** (PLANNED)
- [ ] **Multiple device support** with dynamic device registry
- [ ] **Client implementation** for testing streams
- [ ] **Realistic traffic patterns** (bandwidth simulation, interface state changes)
- [ ] **Configuration management** (admin up/down interfaces via gNMI Set)
- [ ] **Network topology simulation** (device interconnections)

### **Phase 4: Data & Visualization** (PLANNED)
- [ ] **Time-series data storage** (PostgreSQL with proper indexing)
- [ ] **Data aggregation** (5min/1hour/1day rollups)
- [ ] **REST API** for historical data queries
- [ ] **Real-time web dashboard** (React + WebSocket streaming)
- [ ] **Network topology visualization**
- [ ] **Alerting system** (interface down, high utilization)

### **Phase 5: Production Ready** (PLANNED)
- [ ] **Docker containerization** (device simulators, server, dashboard)
- [ ] **Kubernetes deployment** manifests
- [ ] **Auto-scaling** based on device count
- [ ] **Monitoring & logging** (Prometheus metrics, structured logs)
- [ ] **CI/CD pipeline** (GitHub Actions)

## Current Implementation

### **Device Simulator** (`server/device.go`)
- **Concurrent-safe** device with interface management
- **Realistic traffic simulation** (50-400 bytes/update, 5-25 packets/update)
- **Background goroutine** updating counters every 100ms
- **Memory-safe counter reads** with data copying

### **gRPC Protocol** (`proto/telemetry.proto`)
- **InterfaceCounters** message with bytes/packets metrics
- **SubscribeRequest** with device:interface targeting
- **SubscribeResponse** with oneof success/error handling
- **Subscription modes** for different streaming patterns

### **Server Architecture** (`server/main.go`)
- **Multi-device management** with concurrent-safe map
- **RWMutex protection** for high read/low write scenarios
- **Device:interface parsing** for client requests

## Quick Test

```bash
# Run the device simulation
cd server
go run .

# Output: Real-time counter updates
# eth0 - RX: 2164 bytes, TX: 2176 bytes
# eth0 - RX: 4413 bytes, TX: 3996 bytes
# eth0 - RX: 6825 bytes, TX: 6304 bytes
```

## Learning Goals

This project teaches **every technology used in production network monitoring**:

1. **gRPC Streaming** - Handle thousands of concurrent telemetry streams
2. **Concurrent Programming** - Safe multi-threaded network device simulation  
3. **Network Protocols** - Understand gNMI, YANG models, and telemetry patterns
4. **Systems Design** - Build scalable, production-ready network monitoring
5. **Real-time Data** - Stream processing and time-series data management

## Perfect For

- **Network Engineers** learning modern telemetry protocols
- **Backend Engineers** wanting hands-on gRPC and concurrency experience  
- **Systems Programmers** interested in high-performance network applications
- **Anyone** preparing for roles at Cisco, Juniper, or cloud networking companies

## Next Steps

1. **Complete gRPC server** with Subscribe method implementation
2. **Build client** to test streaming functionality  
3. **Add multiple devices** and realistic network topologies
4. **Implement data persistence** and historical querying
5. **Create dashboard** for real-time visualization

---

**Built with ❤️ and lots of ☕ while learning production network monitoring patterns**