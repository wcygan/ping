# Introduction

One of the main reasons to build this project is to see if Buf Schema Registry (BSR, https://buf.build/docs/bsr/introduction) solves
the problem of sharing schemas between services.

## Pushing schema to BSR

### Setup `buf`

Create a [buf.gen.yaml](../buf.gen.yaml) and [buf.yaml](../proto/buf.yaml).

### Make a Schema

Create a proto file, [ping.proto](../proto/ping/v1/ping.proto).

### Setting up BSR 

Visit: https://buf.build/wcygan

Create: https://buf.build/wcygan/ping

Create an auth token: https://buf.build/settings/user

Login:

```bash
buf registry login buf.build
```

### Doing the push

```bash
buf push
```

### Automating the push

Probably you can use github actions to push the schema to BSR. Here is an example of how to do it:

https://github.com/bufbuild/buf-push-action

Yes, I have figured it out in [buf-generate-proto.yaml](../.github/workflows/buf-generate-proto.yaml) github action

## Grand Finale: fetching generated code and running the program

### Fetching generated code

In [server](../server), run these commands to add dependencies to [go.mod](../server/go.mod)

```bash
go get buf.build/gen/go/wcygan/ping/protocolbuffers/go@latest
go get buf.build/gen/go/wcygan/ping/connectrpc/go@latest
go mod tidy
```

Then, I can begin writing code using the imported SDK (see [main.go](../server/main.go)):

```go
import (
	"buf.build/gen/go/wcygan/ping/connectrpc/go/ping/v1/pingv1connect"
	pingv1 "buf.build/gen/go/wcygan/ping/protocolbuffers/go/ping/v1"
	"connectrpc.com/connect"
)

...
```

### Running the program

In [server](../server), run the following command:

```bash
go run .          
2024/10/14 12:20:13 Starting server on :8080
```

Then send a request to the server:

```bash
curl -X POST http://localhost:8080/ping.v1.PingService/Ping \
     -H "Content-Type: application/json" \
     -d '{"timestamp_ms": 1728926331000}'
```