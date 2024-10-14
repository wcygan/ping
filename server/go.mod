module github.com/wcygan/ping/server

go 1.21.6

require (
	buf.build/gen/go/wcygan/ping/connectrpc/go v1.17.0-20241014170349-9ebcb8552d88.1
	buf.build/gen/go/wcygan/ping/protocolbuffers/go v1.35.1-20241014170349-9ebcb8552d88.1
	connectrpc.com/connect v1.17.0
)

require google.golang.org/protobuf v1.35.1 // indirect
