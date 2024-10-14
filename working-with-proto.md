# Introduction

One of the main reasons to build this project is to see if Buf Schema Registry (BSR, https://buf.build/docs/bsr/introduction) solves
the problem of sharing schemas between services.

## Pushing schema to BSR

### Setting it up

Visit: https://buf.build/wcygan

Create: https://buf.build/wcygan/ping

Create an auth token: https://buf.build/settings/user

Login:

```bash
buf registry login buf.build
```

### Doing the push

```bash
github.com/wcygan/ping
```

### Automating the push

Probably you can use github actions to push the schema to BSR. Here is an example of how to do it:

https://github.com/bufbuild/buf-push-action