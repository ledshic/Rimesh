// Rimesh service – P2P virtual LAN daemon.
// Exposes a JSON-RPC API over a local Unix/named-pipe socket that the Tauri
// frontend consumes via Tauri commands / sidecar IPC.
package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ledshic/rimesh/service/internal/network"
)

// RPCRequest is the envelope for incoming JSON-RPC requests.
type RPCRequest struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// RPCResponse is the envelope for JSON-RPC responses.
type RPCResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func main() {
	alias := flag.String("alias", defaultAlias(), "display name shown to peers")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Rimesh service starting (alias=%s)", *alias)

	net, err := network.New(network.Config{
		Alias: *alias,
		OnChange: func() {
			log.Println("[main] peer list changed")
		},
	})
	if err != nil {
		log.Fatalf("create network: %v", err)
	}

	if err := net.Start(); err != nil {
		log.Fatalf("start network: %v", err)
	}

	if err := net.Announce(); err != nil {
		log.Printf("announce: %v", err)
	}

	log.Printf("Virtual IP: %s", net.VirtualIP())

	// Wait for termination signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down…")
	net.Stop()
}

func defaultAlias() string {
	if h, err := os.Hostname(); err == nil {
		return h
	}
	return "rimesh-node"
}
