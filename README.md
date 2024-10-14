# Ping (Client and Server)

This is a basic application which aims to experiment with Buf Schema Registry (BSR) and Cloud Native Postgres (CNPG). 

The application consists of a client and a server. The client sends a request to the server, and the server responds with a message.

## Quickstart

Create a cluster & run the application:

```bash
minikube start
skaffold dev
```

Send a request:

```bash
curl -X POST http://localhost:8080/ping.v1.PingService/Ping \
     -H "Content-Type: application/json" \
     -d '{"timestamp_ms": 1728926331000}'
```

## Tech Stack

- Go
- gRPC
- Protocol Buffers
- Buf
- Buf Schema Registry (push & build on commit?)
- Docker
- Kubernetes
- Minikube (or Kind, or whatever else you have)
- Skaffold
- Postgres
- Cloud Native PostgreSQL - https://cloudnative-pg.io/
- pgx 
- (Optional) NGINX Ingress Controller (or just port forward the server)