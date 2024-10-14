Build the docker image using the following command:

Within [ping]() directory, run the following command:

```bash
docker build -t ping-server -f server/Dockerfile server
```

Note: you can also run it in [server](server) directory like so:

```bash
docker build -t ping-server -f Dockerfile .
```

Run the docker image using the following command:

```bash
docker run -p 8080:8080 ping-server
```