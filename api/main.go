package main

import (
	"fmt"
	"log"

	"github.com/avkos/file-registry/api/config"
	"github.com/avkos/file-registry/api/contracts"
	"github.com/avkos/file-registry/api/handlers"
	"github.com/avkos/file-registry/api/ipfs"
)

// main is the entry point of the application.
func main() {
	if err := config.LoadConfig(); err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}
	fmt.Println("Starting with config:")

	fmt.Printf("Contract: %s\n", config.Config.ContractAddress.Hex())
	fmt.Printf("Ethereum RPC URL: %s\n", config.Config.EthRpcUrl)
	fmt.Printf("IPFS URL: %s\n", config.Config.IpfsUrl)
	fmt.Printf("Port: %s\n", config.Config.Port)

	// Create contract API
	contractAPI, err := contracts.NewContractAPI()
	if err != nil {
		log.Fatalf("Failed to create contract API: %v", err)
	}

	// Create IPFS client
	ipfsClient, err := ipfs.NewIPFSClient(config.Config.IpfsUrl)
	if err != nil {
		log.Fatalf("Failed to create IPFS client: %v", err)
	}

	router := handlers.SetupRouter(contractAPI, ipfsClient)

	addr := ":" + config.Config.Port
	log.Printf("Listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
