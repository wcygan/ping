# Ping (Client and Server)

This is a basic application which aims to experiment with Buf Schema Registry (BSR) and Cloud Native Postgres (CNPG). 

The application consists of a client and a server. The client sends a request to the server, and the server responds with a message.

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