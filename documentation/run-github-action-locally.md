https://github.com/nektos/act

```bash
brew install act
```

We will run [buf-generate-proto.yaml](../.github/workflows/buf-generate-proto.yaml)'s `generate` command:

```bash
act -j generate --container-architecture linux/amd64 -s BUF_TOKEN=<token>
```