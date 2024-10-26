package handler

import (
    "context"
    "log"
    "time"
    "github.com/wcygan/ping/server/service"
    "buf.build/gen/go/wcygan/ping/connectrpc/go/ping/v1/pingv1connect"
    pingv1 "buf.build/gen/go/wcygan/ping/protocolbuffers/go/ping/v1"
    "connectrpc.com/connect"
)

type PingServiceHandler struct {
    service *service.PingService
}

func NewPingServiceHandler(service *service.PingService) *PingServiceHandler {
    return &PingServiceHandler{service: service}
}

func (h *PingServiceHandler) Ping(ctx context.Context, req *connect.Request[pingv1.PingRequest]) (*connect.Response[pingv1.PingResponse], error) {
    timestamp := time.Unix(0, req.Msg.TimestampMs*int64(time.Millisecond)).UTC()
    log.Printf("Received a ping at %s (UTC)", timestamp.Format(time.RFC3339))
    
    if err := h.service.RecordPing(ctx, timestamp); err != nil {
        return nil, err
    }

    return connect.NewResponse(&pingv1.PingResponse{}), nil
}

func (h *PingServiceHandler) PingCount(ctx context.Context, req *connect.Request[pingv1.PingCountRequest]) (*connect.Response[pingv1.PingCountResponse], error) {
    count, err := h.service.GetPingCount(ctx)
    if err != nil {
        return nil, err
    }

    return connect.NewResponse(&pingv1.PingCountResponse{PingCount: count}), nil
}
