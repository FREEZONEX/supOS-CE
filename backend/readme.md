<div align="center">

# Tier0 Edge Backend

[![Go Version](https://img.shields.io/badge/Go-1.24.2+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache_2.0-yellow?style=flat&logo=open-source-initiative)](../LICENSE)
[![Framework](https://img.shields.io/badge/Framework-go--zero-blue?style=flat)](https://github.com/zeromicro/go-zero)

**A high-performance, production-grade backend for industrial IoT data integration**

</div>

---

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [API Documentation](#api-documentation)
- [Configuration](#configuration)
- [Development](#development)
- [Database Schema](#database-schema)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

Tier0 Edge Backend is a robust microservices-based backend system built on the **go-zero framework**. It serves as the core data processing and API layer for the Tier0 IIoT platform, implementing the Unified Namespace (UNS) methodology for industrial data integration.

The backend provides comprehensive REST APIs for device connectivity, data modeling, real-time processing, visualization, and system management.

---

## Architecture

The backend follows a clean architecture pattern with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                     API Gateway Layer                        │
│                  (go-zero REST APIs)                         │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                      Business Logic Layer                    │
│            (Services, Use Cases, Domain Logic)               │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                       Data Access Layer                      │
│         (Repositories, Adapters, External Integrations)      │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    Infrastructure Layer                      │
│    (Databases, Message Queues, External Services)            │
└─────────────────────────────────────────────────────────────┘
```

---

## Features

### Core Capabilities

- **Unified Namespace (UNS) Management**
    - Hierarchical topic modeling
    - Real-time data parsing and validation
    - Tag and label management
    - Tree-based namespace navigation

- **Data Integration & Processing**
    - MQTT message consumption and processing
    - TimescaleDB integration for time-series data
    - PostgreSQL for relational data storage
    - Real-time data streaming and webhook support

- **Device & Protocol Management**
    - Multi-protocol device connectivity
    - Node-RED flow integration
    - Source flow configuration
    - Device status monitoring

- **Visualization & Dashboards**
    - Grafana integration
    - Custom dashboard builder
    - Real-time charting capabilities
    - Template-based visualization

- **Event Processing**
    - Event flow orchestration
    - Rule engine with expression evaluation
    - Alarm management and notifications
    - Trigger functions

- **System Management**
    - User authentication (Keycloak OAuth)
    - Role-based access control
    - Resource management
    - I18n support (multi-language)
    - System configuration

---

## Tech Stack

| Category | Technology |
|----------|-----------|
| **Language** | Go 1.24.2+ |
| **Framework** | [go-zero](https://github.com/zeromicro/go-zero) v1.9.2 |
| **Database ORM** | [GORM](https://gorm.io/) v1.31.0 |
| **Time-Series DB** | TimescaleDB |
| **Relational DB** | PostgreSQL |
| **Message Queue** | MQTT, NATS |
| **Authentication** | Keycloak, JWT |
| **API Gateway** | Kong |
| **Visualization** | Grafana |
| **Flow Editor** | Node-RED |
| **OpenTelemetry** | Jaeger, Prometheus |
| **Cache** | Redis |

---

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go**: 1.24.2 or higher
- **Protocol Buffers**: `protoc`, `protoc-gen-go`, `protoc-gen-grpc-go`
- **Databases**:
    - PostgreSQL 14+
    - TimescaleDB extension
    - Redis (optional, for caching)
- **Message Broker**: MQTT broker (e.g., Mosquitto, EMQX)
- **goctl**: The go-zero CLI tool (fork from [i-Things/go-zero](https://github.com/i-Things/go-zero))

### Installing goctl

```bash
# Clone the go-zero fork
git clone git@github.com:i-Things/go-zero.git

# Install goctl
cd go-zero/tools/goctl
go install

# Verify installation
goctl --version
```

### Installing Protocol Buffer Dependencies

```bash
# One-click install all dependencies
goctl env check --install --verbose --force
```

---

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd Tier0-Edge/backend
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Configure the Application

Edit `etc/backend.yaml` or `etc/backend-local.yaml`:

```yaml
Name: backend
Host: 0.0.0.0
Port: 8080

# Database configuration
Postgres:
  Host: localhost
  Port: 5432
  DBName: tier0
  Username: postgres
  Password: your-password

TimescaleDB:
  Host: localhost
  Port: 5432
  DBName: tier0_ts

# MQTT configuration
Mqtt:
  Broker: tcp://localhost:1883
  ClientID: backend

# Keycloak OAuth
OAuth:
  ClientID: tier0-backend
  ClientSecret: your-secret
  Realm: tier0
  AuthURL: http://localhost:8080/auth
```

### 4. Initialize the Database

```bash
# Run database migrations
go run backend.go migrate
```

### 5. Run the Application

```bash
# Development mode
go run backend.go

# Or build and run
go build -o backend ./backend.go
./backend
```

The API will be available at `http://localhost:8080`

---

## Project Structure

```
backend/
├── backend.go              # Application entry point
├── etc/                    # Configuration files
│   ├── backend.yaml
│   └── backend-local.yaml
├── http/                   # API definitions (go-zero .api files)
│   ├── backend.api        # Main API definition
│   ├── common.api         # Common types
│   └── supos/             # API modules
│       ├── auth.api
│       ├── uns.api
│       ├── system.api
│       ├── devtools.api
│       ├── mount.api
│       ├── sourceflow.api
│       ├── eventflow.api
│       └── ...
├── internal/
│   ├── config/            # Configuration structs
│   ├── handler/           # HTTP request handlers
│   ├── logic/             # Business logic
│   │   └── supos/         # Module-specific logic
│   ├── svc/               # Service context
│   ├── repo/              # Data access layer
│   │   ├── relationDB/    # GORM database models
│   │   └── keycloak/      # Auth repository
│   ├── adapters/          # External service adapters
│   │   ├── grafana/
│   │   ├── kong/
│   │   ├── postgresql/
│   │   ├── timescaledb/
│   │   └── msg_consumer/
│   ├── common/            # Common utilities
│   │   ├── config/
│   │   ├── dto/
│   │   ├── enums/
│   │   ├── cache/
│   │   └── adapter/
│   └── types/             # Type definitions
└── share/                 # Shared utilities
    └── spring/            # Spring-like dependency injection
```

---

## API Documentation

The backend provides RESTful APIs organized by modules:

### Core API Modules

| Module | Description | Base Path |
|--------|-------------|-----------|
| **Auth** | Authentication & Authorization | `/api/auth` |
| **UNS** | Unified Namespace Management | `/api/uns` |
| **System** | System Configuration | `/api/system` |
| **SourceFlow** | Data Source Integration | `/api/sourceflow` |
| **EventFlow** | Event Processing | `/api/eventflow` |
| **Mount** | Resource Mounting | `/api/mount` |
| **DevTools** | Developer Tools | `/api/devtools` |
| **UserManage** | User Management | `/api/user` |
| **Kong** | API Gateway Integration | `/api/kong` |
| **Resource** | Resource Management | `/api/resource` |
| **I18n** | Internationalization | `/api/i18n` |

### Interactive Documentation

Once the server is running, access the interactive API documentation:

- **Swagger UI**: `http://localhost:8080/swagger/`
- **Swagger JSON**: `http://localhost:8080/swagger/swagger.json`

### Example API Calls

```bash
# Get namespace tree
curl -X GET http://localhost:8080/api/uns/namespace/tree \
  -H "Authorization: Bearer <token>"

# Create a new topic
curl -X POST http://localhost:8080/api/uns/namespace/topic \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "temperature",
    "path": "factory/line1/equipment1",
    "dataType": "number",
    "unit": "°C"
  }'

# Query time-series data
curl -X POST http://localhost:8080/api/uns/data/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "topic": "factory/line1/equipment1/temperature",
    "startTime": "2024-01-01T00:00:00Z",
    "endTime": "2024-01-01T23:59:59Z",
    "aggregation": "avg",
    "interval": "1h"
  }'
```

---

## Configuration

### Environment Variables

Key environment variables (typically set in `etc/backend.yaml`):

| Variable | Description | Default |
|----------|-------------|---------|
| `Host` | Server host | `0.0.0.0` |
| `Port` | Server port | `8080` |
| `Log.Mode` | Log mode (console/file) | `console` |
| `Log.Level` | Log level (info/error/debug) | `info` |
| `Postgres.Host` | PostgreSQL host | `localhost` |
| `Postgres.Port` | PostgreSQL port | `5432` |
| `Postgres.DBName` | Database name | `tier0` |
| `Mqtt.Broker` | MQTT broker URL | `tcp://localhost:1883` |
| `OAuth.ClientID` | OAuth client ID | - |
| `OAuth.ClientSecret` | OAuth client secret | - |

---

## Development

### Adding a New RPC Service

```bash
# Create new RPC service
goctl rpc new <service-name> --style=goZero -m
```

### Adding a New API Service

```bash
# Create new API service
goctl api new <service-name> --style=goZero
```

### Adding a New Database Table

1. Generate GORM models using [sql2gorm](https://sql2gorm.mccode.info/)
2. Place models in `internal/repo/relationDB/model.go`
3. Copy `internal/repo/relationDB/example.go` and rename to your table
4. Replace `Example` with your table name
5. Implement custom functions as needed

### Generating API Code from Definitions

```bash
# Generate handler, logic, and types from API definition
goctl api go -api http/backend.api -dir ./ --style=goZero

# Generate Swagger documentation
goctl api swagger -filename swagger.json -api http/backend.api -dir ./http

# Generate access control code
goctl api access -api http/backend.api -dir ./http
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/repo/relationDB
```

---

## Database Schema

The backend uses two databases:

### PostgreSQL (Relational Data)

Stores system configuration, user data, namespace definitions, and metadata.

Key tables:
- `uns_namespaces` - Namespace definitions
- `uns_tags` - Tag configurations
- `uns_labels` - Label definitions
- `uns_alarms` - Alarm configurations
- `dashboards` - Dashboard definitions
- `nodered_flows` - Node-RED flow definitions
- `sys_users` - User information
- `sys_roles` - Role definitions

### TimescaleDB (Time-Series Data)

Stores high-volume time-series data from devices.

Key hypertables:
- `uns_data_points` - Time-series data points
- `alarm_history` - Alarm history records

---

## Deployment

### Building for Production

```bash
# Build for Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o backend ./backend.go

# Build for macOS
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o backend ./backend.go

# Build for Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o backend.exe ./backend.go
```

### Docker Deployment

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o backend ./backend.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/backend .
COPY --from=builder /app/etc ./etc
CMD ["./backend", "-f", "etc/backend.yaml"]
```

### Deployment Script

The project includes a deployment script for remote deployment:

```bash
# Build, package, and deploy to remote server
cd /home/supOS-V1.2.0.0-M-25102214-T5/
sh rebuild_start_backend.sh
```

---

## Development Environment Setup

For local development with the testing environment:

1. **Server**: `supos-community-edge`
2. **Access**: `http://100.100.100.22:34098/home`
3. **Database**: `100.100.100.22:34099`

### Deploying to Development Environment

```bash
# Build executable
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o backend ./backend.go

# Upload to server
scp backend user@100.100.100.22:/home/supos-V1.2.0.0-M-25102214-T5/gobackend/

# One-key deployment script (build image and restart service)
cd /home/supOS-V1.2.0.0-M-25102214-T5/
sh rebuild_start_backend.sh
```

---

## Architecture Patterns

### Clean Architecture

The backend follows clean architecture principles:

1. **Domain Layer**: Core business logic and entities
2. **Application Layer**: Use cases and application services
3. **Infrastructure Layer**: External integrations and data access
4. **Interface Layer**: HTTP handlers and API endpoints

### Dependency Injection

Uses a Spring-like dependency injection framework (`share/spring`):

```go
// Register beans
spring.RegisterBean[*svc.ServiceContext](ctx)
spring.RefreshBeanContext()

// Publish events
spring.PublishEvent(&event.ContextRefreshedEvent{SvcContext: ctx})
```

### Repository Pattern

Data access abstracted through repositories:

```go
type UnsNamespaceRepo interface {
    Insert(ctx context.Context, data *UnsNamespace) error
    FindOne(ctx context.Context, id int64) (*UnsNamespace, error)
    Update(ctx context.Context, data *UnsNamespace) error
    Delete(ctx context.Context, id int64) error
}
```

---

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](../CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Development Guidelines

- Follow the existing code style and structure
- Write unit tests for new features
- Update API documentation (`.api` files)
- Ensure all tests pass before submitting
- Use meaningful commit messages

---

## Troubleshooting

### Common Issues

**Issue**: `goctl: command not found`
```bash
# Solution: Install goctl from the correct fork
git clone git@github.com:i-Things/go-zero.git
cd go-zero/tools/goctl && go install
```

**Issue**: Database connection errors
```bash
# Solution: Verify PostgreSQL/TimescaleDB is running
psql -h localhost -U postgres -d tier0

# Check connection in etc/backend.yaml
```

**Issue**: MQTT connection failures
```bash
# Solution: Verify MQTT broker is accessible
mosquitto_sub -h localhost -t "test/topic" -v
```

---

## Performance

- **API Response Time**: < 100ms (p95)
- **Throughput**: 10,000+ requests/second
- **Database Query Optimization**: Indexed queries, connection pooling
- **Caching**: Redis caching for frequently accessed data

---

## Monitoring & Observability

The backend integrates with several monitoring tools:

- **OpenTelemetry**: Distributed tracing
- **Jaeger**: Trace visualization
- **Prometheus**: Metrics collection
- **pprof**: Profiling (built into go-zero)

---

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](../LICENSE) file for details.

---

## Support & Contact

- 📖 [Documentation](https://tier0edge.vercel.app/)
- 🐛 [Issue Tracker](https://github.com/freezonex/Tier0-Edge/issues)
- 💬 [Discussions](https://github.com/freezonex/Tier0-Edge/discussions)

---

## Acknowledgments

Built with:

- [go-zero](https://github.com/zeromicro/go-zero) - Go microservices framework
- [GORM](https://gorm.io/) - ORM library
- [TimescaleDB](https://www.timescale.com/) - Time-series database
- [Grafana](https://grafana.com/) - Visualization platform
- [Node-RED](https://nodered.org/) - Flow-based programming
- All other open-source contributors

---

<div align="center">

**Built with ❤️ for the Industrial IoT Community**

</div>
