package main

import (
	"buf.build/gen/go/wcygan/ping/connectrpc/go/ping/v1/pingv1connect"
	pingv1 "buf.build/gen/go/wcygan/ping/protocolbuffers/go/ping/v1"
	"connectrpc.com/connect"
	"context"
	"log"
	"net/http"
	"time"
)

// PingServiceServer implements the PingService interface.
type PingServiceServer struct{}

// Ping handles the Ping RPC.
func (s *PingServiceServer) Ping(ctx context.Context, req *connect.Request[pingv1.PingRequest]) (*connect.Response[pingv1.PingResponse], error) {
	timestamp := time.Unix(0, req.Msg.TimestampMs*int64(time.Millisecond)).UTC()
	log.Printf("Received ping at %s (UTC)", timestamp.Format(time.RFC3339))
	return connect.NewResponse(&pingv1.PingResponse{}), nil
}

// PingCount handles the PingCount RPC.
func (s *PingServiceServer) PingCount(ctx context.Context, req *connect.Request[pingv1.PingCountRequest]) (*connect.Response[pingv1.PingCountResponse], error) {
	return connect.NewResponse(&pingv1.PingCountResponse{PingCount: 42}), nil
}

func main() {
	mux := http.NewServeMux()
	server := &PingServiceServer{}

	// Register the PingService with the Connect server.
	mux.Handle(pingv1connect.NewPingServiceHandler(server))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
