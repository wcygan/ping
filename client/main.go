package main

import (
	"buf.build/gen/go/wcygan/ping/connectrpc/go/ping/v1/pingv1connect"
	pingv1 "buf.build/gen/go/wcygan/ping/protocolbuffers/go/ping/v1"
	"connectrpc.com/connect"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger *zap.Logger

var rootCmd = &cobra.Command{
	Use:   "pingcli",
	Short: "Ping CLI is a tool to interact with the Ping server",
}

var address string

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Send a ping request to the server",
	Run: func(cmd *cobra.Command, args []string) {
		client := pingv1connect.NewPingServiceClient(http.DefaultClient, address)
		pingRequest := &pingv1.PingRequest{
			TimestampMs: time.Now().UnixMilli(),
		}
		resp, err := client.Ping(context.Background(), connect.NewRequest(pingRequest))
		if err != nil {
			logger.Fatal("Ping request failed",
				zap.Error(err),
				zap.Int64("timestamp_ms", pingRequest.TimestampMs))
		}

		logger.Info("Successfully pinged",
			zap.Time("timestamp", time.UnixMilli(pingRequest.TimestampMs).UTC()),
			zap.Any("response", resp.Msg))
	},
}

var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Send a ping count request to the server",
	Run: func(cmd *cobra.Command, args []string) {
		client := pingv1connect.NewPingServiceClient(http.DefaultClient, address)
		pingCountRequest := &pingv1.PingCountRequest{}
		resp, err := client.PingCount(context.Background(), connect.NewRequest(pingCountRequest))
		if err != nil {
			logger.Fatal("PingCount request failed", zap.Error(err))
		}
		logger.Info("Received ping count", 
			zap.Int64("count", resp.Msg.PingCount))
	},
}

func init() {
	// Initialize logger
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	// Setup command flags
	rootCmd.PersistentFlags().StringVar(&address, "address", "http://localhost:8080", "server address")
	rootCmd.AddCommand(pingCmd)
	rootCmd.AddCommand(countCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
