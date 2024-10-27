module github.com/wcygan/ping/client

go 1.22.0

require (
	buf.build/gen/go/wcygan/ping/connectrpc/go v1.17.0-20241014170349-9ebcb8552d88.1
	buf.build/gen/go/wcygan/ping/protocolbuffers/go v1.35.1-20241014170349-9ebcb8552d88.1
	connectrpc.com/connect v1.17.0
	github.com/spf13/cobra v1.8.1
	go.uber.org/zap v1.27.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)
