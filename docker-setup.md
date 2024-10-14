Build the docker image using the following command:

Run these commands within [ping]() directory.

Build the container images:

```bash
docker build -t ping-server -f server/Dockerfile server
docker build -t ping-client -f client/Dockerfile client
```

Create a network:

```bash
docker network create ping-network
```

Run the server container:

```bash
docker run -p 8080:8080 --network ping-network --name ping-server ping-server
```

Run the client container:

```bash
docker run --network ping-network ping-client ping --address http://ping-server:8080
```

You can also run the binary directly in [client](client) folder:

```bash
go run . ping --address http://:8080
```