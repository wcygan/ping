# ConnectRPC, Buf, and Kubernetes Setup Guide

## Prerequisites

- Go 1.22+
- Docker
- kubectl
- skaffold
- buf CLI

## 1. Protocol Buffer Definition
The service is defined in `proto/ping/v1/ping.proto`:
```protobuf
service PingService {
  rpc Ping (PingRequest) returns (PingResponse);
  rpc PingCount (PingCountRequest) returns (PingCountResponse);
}
```

## 2. Buf Configuration
### buf.yaml
```yaml
version: v1
name: buf.build/wcygan/ping
breaking:
  use:
    - FILE
lint:
  use:
    - DEFAULT
```

### buf.gen.yaml
```yaml
version: v2
managed:
  enabled: true
  override:
      - file_option: go_package_prefix
        value: github.com/wcygan/proto/generated/go
plugins:
  - remote: buf.build/protocolbuffers/go
    out: gen
    opt: paths=source_relative
  - remote: buf.build/connectrpc/go
    out: gen
    opt: paths=source_relative
```

## 3. GitHub Actions Integration
The `.github/workflows/buf-generate-proto.yaml` workflow automatically pushes schema updates to the Buf Schema Registry:
```yaml
name: Generate Go code from Buf schema
on:
  push:
    branches:
      - main
jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: bufbuild/buf-setup-action@v1.45.0
      - uses: bufbuild/buf-push-action@v1.2.0
        with:
          buf_token: ${{ secrets.BUF_TOKEN }}
```

## 4. ConnectRPC Server Implementation
The server implementation uses the generated code to handle RPC calls:
```go
func (h *PingServiceHandler) Ping(ctx context.Context, req *connect.Request[pingv1.PingRequest]) (*connect.Response[pingv1.PingResponse], error) {
    timestamp := time.Unix(0, req.Msg.TimestampMs*int64(time.Millisecond)).UTC()
    if err := h.service.RecordPing(ctx, timestamp); err != nil {
        return nil, err
    }
    return connect.NewResponse(&pingv1.PingResponse{}), nil
}
```

## 5. Kubernetes Deployment
The service is deployed to Kubernetes using:

### Dockerfile
```dockerfile
FROM golang:1.22 as builder
WORKDIR /app/server
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./main.go

FROM alpine:3.14
COPY --from=builder /app/server/server /app/server
COPY --from=builder /app/server/migrations /app/migrations
ENTRYPOINT ["/app/server"]
```

### Kubernetes Deployment (server.yaml)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: server
spec:
  replicas: 1
  template:
    spec:
      containers:
        - name: server
          image: server
          ports:
            - containerPort: 8080
          env:
            - name: DB_HOST
              value: pg-cluster-rw
```

## Key Features
- **Buf Schema Registry**: Centralized schema management
- **ConnectRPC**: Modern RPC framework with HTTP/1.1 and HTTP/2 support
- **Kubernetes**: Container orchestration with:
  - Health checks
  - Environment configuration
  - Database migration handling
  - Service discovery
  - Horizontal Pod Autoscaling

## Service Dependencies

### PostgreSQL
- Managed by CloudNativePG operator
- Automatic failover and high availability
- Credentials managed via Kubernetes secrets

### Kafka
- Used for event streaming
- Topics:
  - ping-events: Records all ping events
- Managed by Strimzi operator

## Local Development Setup

1. Start local dependencies:
```bash
# Start PostgreSQL
docker run -d --name postgres \
  -e POSTGRES_USER=pinguser \
  -e POSTGRES_PASSWORD=ping123 \
  -e POSTGRES_DB=pingdb \
  -p 5432:5432 \
  postgres:15

# Start Kafka
docker run -d --name kafka \
  -p 9092:9092 \
  -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
  bitnami/kafka:latest
```

2. Run database migrations:
```bash
cd server
go run main.go -migrate-only
```

3. Start the server:
```bash
go run main.go
```

## Testing Strategy

1. **Unit Tests**
   - Repository layer: `repository/ping_repository_test.go`
   - Service layer: `service/ping_service_test.go`
   - Handler layer: `handler/ping_handler_test.go`

2. **Integration Tests**
   - Database integration tests with testcontainers
   - Kafka integration tests with embedded broker

3. **End-to-End Tests**
   - Full API tests using real ConnectRPC clients
   - Kubernetes deployment tests using kind

## Development Workflow
1. Update proto definitions
2. Generate code: `buf generate`
3. Run tests: `go test ./...`
4. Push to main branch
5. GitHub Actions automatically updates Buf Schema Registry
6. Build and deploy server using Skaffold

## Monitoring and Observability

- Health checks via `/healthz` endpoint
- Structured logging with Zap
- Metrics exposed for Prometheus
- Distributed tracing with OpenTelemetry
