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

| Technology                                                                                                                                                              | Purpose                                                               |
|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------|
| [Go](https://go.dev/)                                                                                                                                                   | Programming language used for the application                         |
| [ConnectRPC](https://connectrpc.com/)                                                                                                                                   | Framework for remote procedure calls                                  |
| [Protocol Buffers](https://protobuf.dev/)                                                                                                                               | Serialization format used with ConnectRPC                             |
| [Buf](https://buf.build/)                                                                                                                                               | Tool for working with Protocol Buffers                                |
| [Buf Schema Registry](https://buf.build/product/bsr)                                                                                                                    | Registry for storing and managing Protocol Buffers schemas            |
| [Docker](https://www.docker.com/)                                                                                                                                       | Containerization platform                                             |
| [Kubernetes](https://kubernetes.io/)                                                                                                                                    | Container orchestration platform                                      |
| [Minikube](https://minikube.sigs.k8s.io/docs/start/?arch=%2Fmacos%2Farm64%2Fstable%2Fbinary+download) (or [Kind](https://kind.sigs.k8s.io/), or whatever else you have) | Tool for running a local Kubernetes cluster                           |
| [Skaffold](https://skaffold.dev/)                                                                                                                                       | Tool for building and deploying applications on Kubernetes            |
| [Postgres](https://www.postgresql.org/)                                                                                                                                 | Relational database management system                                 |
| [CloudNativePG](https://cloudnative-pg.io/)                                                                                                                             | Operator for managing PostgreSQL on Kubernetes                        |
| [pgx](https://github.com/jackc/pgx)                                                                                                                                     | PostgreSQL driver and toolkit for Go                                  |
| [Orbstack](https://orbstack.dev/) (or [Docker Desktop](https://www.docker.com/products/docker-desktop/))                                                                | Virtualized environment for running containers                        |
