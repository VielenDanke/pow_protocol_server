package main

import (
	"context"
	"github.com/vielendanke/pow_protocol_server/server"
	"log"
	"os"
)

func main() {
	networkType := os.Getenv("SERVER_NETWORK_TYPE")
	serverAddress := os.Getenv("SERVER_ADDRESS")
	ds, dsErr := server.NewDefaultServer(networkType, serverAddress)
	if dsErr != nil {
		log.Printf("ERROR: cannot start server due to error: %s\n", dsErr)
		return
	}
	log.Printf("INFO: Server stopped: %s", ds.Start(context.Background()))
}
