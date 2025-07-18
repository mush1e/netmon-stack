# NetMon Stack - Real-Time Network Monitoring System

A **gRPC-based network telemetry learning project** that simulates Cisco network devices and provides real-time interface counter streaming. Built to learn and demonstrate modern network monitoring technologies.

## Project Overview

**NetMon Stack** is a well-structured learning project that demonstrates:
- **gRPC & Protocol Buffers** for high-performance network communication
- **gNMI-style telemetry** for real-time network data streaming  
- **Go concurrency patterns** with goroutines, channels, and mutexes
- **Realistic network simulation** with traffic patterns and interface counters
- **Production-ready architecture** with proper error handling and thread safety
- **Message queue integration** with Redis pub/sub for scalable event streaming
- **Containerization & orchestration** with Docker and Kubernetes

## Architecture

### **Current Implementation**
```
┌─────────────────┐    gRPC Stream    ┌─────────────────┐
│                 │◄─────────────────►│                 │
│   grpcurl       │   Subscribe()     │  gRPC Server    │
│  (Test Client)  │   device:iface    │                 │
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

### **Target Architecture (Event-Driven)**
```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│ Device Sim 1    │────▶│                 │◄───▶│ Web Dashboard   │
├─────────────────┤     │                 │     ├─────────────────┤
│ Device Sim 2    │────▶│  Redis Pub/Sub  │◄───▶│ Go Client App   │
├─────────────────┤     │                 │     ├─────────────────┤
│ Device Sim 3    │────▶│   (Message      │◄───▶│ Alert Service   │
└─────────────────┘     │    Queue)       │     ├─────────────────┤
                        └─────────────────┘     │ Database Writer │
                                                └─────────────────┘
```

### **Containerized Production Deployment**
```
┌────────────────────────────────────────────────────────────┐
│                   Kubernetes Cluster                       │   
│                                                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Router1    │  │   Switch1    │  │   Switch2    │      │
│  │  Simulator   │  │  Simulator   │  │  Simulator   │      │
│  │   (Pod)      │  │   (Pod)      │  │   (Pod)      │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│         │                  │                  │            │
│         └──────────────────┼──────────────────┘            │
│                            │                               │
│                    ┌──────────────┐                        │
│                    │ gRPC Server  │                        │
│                    │   (Pod)      │                        │
│                    └──────────────┘                        │
│                            │                               │
│                    ┌──────────────┐                        │
│                    │ Redis Queue  │                        │
│                    │   (Pod)      │                        │
│                    └──────────────┘                        │
│                            │                               │
│                    ┌──────────────┐                        │
│                    │  Dashboard   │                        │
│                    │   (Pod)      │                        │
│                    └──────────────┘                        │
└────────────────────────────────────────────────────────────┘
```

## Features & Roadmap

### **✅ Phase 1: Protocol Foundation** (COMPLETED)
- [x] **Protocol Buffer definitions** for network telemetry data
- [x] **gNMI-style gRPC service** with Subscribe streaming
- [x] **Device simulation** with realistic interface counters
- [x] **Concurrent counter updates** using goroutines (100ms intervals)
- [x] **Thread-safe operations** with RWMutex patterns

### **✅ Phase 2: gRPC Server Implementation** (COMPLETED)
- [x] **Complete server implementation** with device management
- [x] **Subscribe method** with STREAM and ONCE modes
- [x] **String parsing** for "device:interface" format (router1:eth0)
- [x] **Enterprise error handling** with custom ErrorCode enum
- [x] **Multiple subscription modes** (STREAM, ONCE, POLL stub)
- [x] **Production-grade streaming** with context cancellation
- [x] **gRPC reflection** for service discovery

### **🚧 Phase 3: Client Applications** (IN PROGRESS)
- [ ] **Go gRPC client** for programmatic access
- [ ] **Command-line interface** for network operators
- [ ] **Multiple device support** with dynamic device registry
- [ ] **Redis pub/sub integration** for event-driven architecture
- [ ] **Realistic traffic patterns** (bandwidth simulation, interface state changes)

### **📋 Phase 4: Data & Visualization** (PLANNED)
- [ ] **Time-series data storage** (PostgreSQL/InfluxDB with proper indexing)
- [ ] **Data aggregation** (5min/1hour/1day rollups)
- [ ] **REST API** for historical data queries
- [ ] **Real-time web dashboard** (React + WebSocket streaming)
- [ ] **Network topology visualization** with D3.js
- [ ] **Alerting system** (interface down, high utilization thresholds)

### **📋 Phase 5: Production Ready** (PLANNED)
- [ ] **Docker containerization** (multi-stage builds, optimized images)
- [ ] **Kubernetes deployment** manifests with auto-scaling
- [ ] **Health checks** and rolling updates
- [ ] **Monitoring & logging** (Prometheus metrics, structured logs)
- [ ] **CI/CD pipeline** (GitHub Actions, automated testing)
- [ ] **Load testing** and performance optimization

## Current Implementation

### **Device Simulator** (`server/device.go`)
- **Thread-safe** device with interface management
- **Realistic traffic simulation** (50-400 bytes/update, 5-25 packets/update)
- **Background goroutines** updating counters every 100ms
- **Safe counter reads** with data copying to prevent race conditions
- **Good concurrency practices** using RWMutex for read/write scenarios

### **gRPC Protocol** (`proto/telemetry.proto`)
- **InterfaceCounters** message with bytes/packets metrics and timestamps
- **SubscribeRequest** with device:interface targeting and configurable intervals
- **SubscribeResponse** with oneof success/error handling
- **Custom ErrorCode enum** (DOES_NOT_EXIST, PERMISSION_DENIED, NOT_ACTIVE)
- **Subscription modes** for different streaming patterns

### **gRPC Server** (`server/main.go`)
- **Well-structured implementation** of NetworkTelemetryServer interface
- **Multi-device management** with thread-safe device map
- **Context-aware streaming** with graceful client disconnection handling
- **Good error handling** with proper gRPC status codes
- **Input validation** for device:interface format
- **Configurable streaming intervals** with sensible defaults

## Quick Start & Testing

### **Start the Server**
```bash
cd server
go run .
# Output: Starting gRPC server on :50051
```

### **Test with grpcurl**
```bash
# Install grpcurl
brew install grpcurl

# List available services
grpcurl -plaintext localhost:50051 list

# Test ONCE mode (single response)
grpcurl -plaintext -d '{
  "interface_name": "router1:eth0",
  "mode": "ONCE"
}' localhost:50051 NetworkTelemetry/Subscribe

# Test STREAM mode (continuous updates)
grpcurl -plaintext -d '{
  "interface_name": "router1:eth0", 
  "interval_ms": 2000,
  "mode": "STREAM"
}' localhost:50051 NetworkTelemetry/Subscribe
```

### **Expected Output (Live Streaming)**
```json
{
  "counters": {
    "interfaceName": "eth0",
    "bytesRx": "156814",
    "bytesTx": "156079", 
    "packetsRx": "10381",
    "packetsTx": "10388",
    "timestamp": 347818225
  },
  "responseTimestamp": "1752694475058"
}
```

## Learning Outcomes

This project demonstrates **solid understanding** of:

### **Backend Engineering**
- **Go concurrency fundamentals** - goroutines, channels, mutexes, context cancellation
- **gRPC streaming basics** - server streaming, client lifecycle management
- **Protocol Buffers** - schema design, code generation, serialization
- **Error handling patterns** - custom error types, graceful error responses

### **Systems Programming** 
- **Thread-safe programming** - race condition prevention, safe data access
- **Network programming basics** - TCP listeners, HTTP/2 concepts
- **Resource management** - goroutine lifecycle, basic memory management

### **Network Engineering Concepts**
- **gNMI protocol basics** - telemetry streaming, subscription modes
- **Network device simulation** - interface modeling, counter patterns
- **Telemetry data structures** - counters, timestamps, device hierarchies

### **Modern Development Practices**
- **Code organization** - clean structure, separation of concerns
- **Testing approaches** - using tools like grpcurl for validation
- **Documentation** - clear README, code comments
- **Version control** - proper Git practices

## Technical Stack

- **Go 1.23** - Systems programming with excellent concurrency primitives
- **gRPC** - High-performance RPC framework with HTTP/2 multiplexing
- **Protocol Buffers** - Efficient binary serialization and schema evolution
- **Redis** - In-memory message queuing and pub/sub (planned)
- **Docker** - Containerization and deployment (planned)
- **Kubernetes** - Container orchestration and auto-scaling (planned)

## Next Steps

1. **Build Go gRPC client** for programmatic testing and integration
2. **Add Redis pub/sub** for event-driven architecture and scalability
3. **Create web dashboard** with real-time charts and network topology
4. **Implement time-series storage** for historical data analysis
5. **Containerize with Docker** and deploy on Kubernetes

---

**Built as a learning project to understand modern network monitoring and distributed systems concepts**
