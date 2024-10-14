# Introduction

One of the main reasons to build this project is to see if Buf Schema Registry (BSR, https://buf.build/docs/bsr/introduction) solves
the problem of sharing schemas between services.

## Pushing schema to BSR

### Setup `buf`

Create a `buf.gen.yaml` file with the following content:

(File is located at: `buf.gen.yaml`)

```yaml
version: v1
managed:
  enabled: true
  go_package_prefix:
    default: github.com/wcygan/proto/generated/go
plugins:
  - name: go
    out: generated/go
    opt: paths=source_relative
```

Create a `buf.yaml` file with the following content:

(File is located at: `proto/buf.yaml`)

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

### Make a Schema

(File is located at: `proto/ping/v1/ping.proto`)

Create a file called `ping.proto` with the following content:

```proto
syntax = "proto3";

package ping.v1;

service PingService {
  rpc Ping (PingRequest) returns (PingResponse);
  rpc PingCount (PingCountRequest) returns (PingCountResponse);
}

message PingRequest {
  int64 timestamp_ms = 1;
}

message PingResponse {}

message TimeRange {
  int64 start = 1;
  int64 end = 2;
}

message PingCountRequest {
  optional TimeRange time_range = 1;
}

message PingCountResponse {
  int64 ping_count = 1;
}
```

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