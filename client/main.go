package main

import (
	"buf.build/gen/go/wcygan/ping/connectrpc/go/ping/v1/pingv1connect"
	pingv1 "buf.build/gen/go/wcygan/ping/protocolbuffers/go/ping/v1"
	"connectrpc.com/connect"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

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
		_, err := client.Ping(context.Background(), connect.NewRequest(pingRequest))
		if err != nil {
			log.Fatalf("Ping request failed: %v", err)
		}

		log.Printf("Successfully pinged at %s (UTC)", time.UnixMilli(pingRequest.TimestampMs).Format(time.RFC3339))
	},
}

var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Send a ping count request to the server",
	Run: func(cmd *cobra.Command, args []string) {
		client := pingv1connect.NewPingServiceClient(http.DefaultClient, address)
		pingCountRequest := &pingv1.PingCountRequest{}
		pingCountResponse, err := client.PingCount(context.Background(), connect.NewRequest(pingCountRequest))
		if err != nil {
			log.Fatalf("PingCount request failed: %v", err)
		}
		fmt.Printf("PingCount response: %v\n", pingCountResponse.Msg.PingCount)
	},
}

func init() {
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
