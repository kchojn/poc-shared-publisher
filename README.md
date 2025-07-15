# POC Shared Publisher - Phase 1

A proof-of-concept implementation of a shared publisher for cross-chain transaction coordination. This is **Phase 1** -
a simple message relay system between rollup sequencers.

## Architecture

The system consists of three main components:

```
┌─────────────┐    XTRequest     ┌─────────────────┐    broadcast    ┌─────────────┐
│             │ ──────────────>  │     Shared      │ ──────────────>  │             │
│ Sequencer A │                  │    Publisher    │                  │ Sequencer B │
│             │                  │     (Hub)       │                  │             │
└─────────────┘                  └─────────────────┘                  └─────────────┘
```

### Components

1. **Shared Publisher (SP)** - Central hub that receives and broadcasts transactions
    - Listens on port `8080` for TCP connections from sequencers
    - Exposes metrics and health endpoints on port `8081`
    - Maintains persistent connections with multiple sequencers

2. **Sequencer A** - Sends cross-chain transaction requests
    - Connects to Shared Publisher via TCP
    - Sends `XTRequest` messages containing transaction data

3. **Sequencer B** - Receives broadcasted transactions
    - Connects to Shared Publisher via TCP
    - Receives broadcasted `XTRequest` messages from other sequencers

## Message Flow (Phase 1)

1. **Connection Setup**: Sequencers establish TCP connections to the Shared Publisher
2. **Transaction Submission**: Sequencer A sends an `XTRequest` message to the SP
3. **Message Broadcast**: SP receives the message and broadcasts it to all other connected sequencers
4. **Transaction Processing**: Sequencer B receives the broadcasted message and can process it

### Message Types

Based on the protobuf definition in `api/proto/messages.proto`:

```protobuf
// User request
message XTRequest {
  repeated TransactionRequest transactions = 1;
}

message TransactionRequest {
  bytes chainID = 1;
  repeated bytes transaction = 2;
}
```

## Quick Start

### Prerequisites

- Go 1.24+
- Docker and Docker Compose
- Make

### Running with Docker

```bash
# Build and run the system
make docker-run

# Or manually
docker-compose up --build
```

### Running Locally

```bash
# Build the application
make build

# Run the publisher
make run

# Or directly
./bin/poc-shared-publisher -config configs/config.yaml
```

### Testing the System

Use the provided Python scripts to simulate sequencers:

```bash
# Terminal 1: Start the publisher
make docker-run

# Terminal 2: Send a test transaction
python3 scripts/send_request.py

# Terminal 3: Run multiple clients simulation
python3 scripts/multiple_clients.py
```

## Configuration

The system uses a YAML configuration file (`configs/config.yaml`):

```yaml
server:
  listen_addr: ":8080"          # TCP port for sequencer connections
  read_timeout: 30s             # Connection read timeout
  write_timeout: 30s            # Connection write timeout
  max_message_size: 10485760    # 10MB max message size
  max_connections: 10           # Max concurrent connections (Phase 1)

metrics:
  enabled: true                 # Enable Prometheus metrics
  port: 8081                    # HTTP port for metrics

log:
  level: info                   # Log level
  pretty: false                 # JSON logging for Loki
  output: stdout                # Output to stdout (Loki integration)
```

### Environment Variables

All configuration values can be overridden using environment variables:

```bash
# Server configuration
export SERVER_LISTEN_ADDR=":9090"
export SERVER_MAX_CONNECTIONS=200
export SERVER_READ_TIMEOUT=60s

# Metrics configuration
export METRICS_PORT=3000
export METRICS_ENABLED=false

# Logging configuration
export LOG_LEVEL=debug
export LOG_PRETTY=true
```

**Pattern**: `<SECTION>_<KEY>` (dots replaced with underscores, all uppercase)

**Examples**:

- `server.listen_addr` → `SERVER_LISTEN_ADDR`
- `metrics.port` → `METRICS_PORT`
- `log.level` → `LOG_LEVEL`

**Priority**: ENV variables > YAML config > Default values

## Monitoring

### Metrics

Prometheus metrics are exposed on `http://localhost:8081/metrics`:

- `crosschain_transactions_total` - Total cross-chain transactions processed
- `connections_active` - Number of active sequencer connections
- `broadcasts_total` - Total messages broadcasted
- `message_processing_duration_seconds` - Message processing time

### Health Checks

- **Health**: `http://localhost:8081/health` - System health status
- **Ready**: `http://localhost:8081/ready` - Readiness status (has connections)
- **Stats**: `http://localhost:8081/stats` - Publisher statistics
- **Connections**: `http://localhost:8081/connections` - Active connections info

### Prometheus Setup

Use the provided Prometheus configuration:

```bash
# The config is optimized for Phase 1 development
cat monitoring/prometheus/prometheus.yml
```

## Development

### Building

```bash
# Build binary
make build

# Run tests
make test

# Run tests with coverage
make coverage

# Run linters
make lint

# Generate protobuf files
make proto
```

### Project Structure

```
├── api/proto/              # Protobuf definitions
├── cmd/publisher/          # Main application entry point
├── configs/               # Configuration files
├── internal/
│   ├── config/           # Configuration management
│   ├── network/          # TCP server/client implementation
│   ├── proto/            # Generated protobuf files
│   └── publisher/        # Core publisher logic
├── monitoring/           # Prometheus configuration
├── pkg/
│   ├── logger/          # Logging utilities
│   └── metrics/         # Prometheus metrics
└── scripts/             # Development and testing scripts
```

### Protocol

The system uses a simple TCP protocol with length-prefixed protobuf messages:

```
[4 bytes length][protobuf message]
```
