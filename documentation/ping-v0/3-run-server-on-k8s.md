## 0: Orbstack (or Docker Desktop)

The first thing you need is an environment to run virtualized containers.

Orbstack, https://orbstack.dev/, is suitable for this purpose. I use this because it is faster than Docker Desktop on my Macbook Pro M1 Max.

Alternatively, you can use Docker Desktop to run a Kubernetes cluster locally (https://www.docker.com/products/docker-desktop/).

I typically open up the Orbstack (or Docker Desktop) program and then proceed to start my MiniKube cluster.

## 1: Minikube

Minikube creates a single-node Kubernetes cluster on your local machine. It is suitable for development and testing purposes.

Visit https://minikube.sigs.k8s.io/docs/start/ for getting started.

With your virtualized environment running, you can start Minikube with the following command:

```bash
minikube start
```

## 2: Skaffold

With these prerequisites in place, you can now use Skaffold to build and deploy your application to the Kubernetes cluster.

To understand how this application is deployed, reference [skaffold.yaml](../skaffold.yaml) and run through all of
the configurations that are specified.

The application is set up to have a repeatable build and deploy process.

## 3. Test the application

Visit the [Quickstart](../README.md#quickstart) section of the README to see how to run and test the application.