package main

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	blockCheckInterval = 10 * time.Second
)

func checkBlock(rpcUrl string) (bool, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Get the latest block number
	latestBlock, err := getLatestBlock(client)
	if err != nil {
		log.Fatal(err)
	}
	if !checkService() {
		return false, nil
	}
	time.Sleep(blockCheckInterval)
	latestBlockAgain, err := getLatestBlock(client)
	if err != nil {
		log.Fatal(err)
	}

	// If the latest block has changed, log the error and exit the loop
	if latestBlockAgain.NumberU64() == latestBlock.NumberU64() {
		return false, nil
	}
	return true, nil
}

// Get the latest block from the connected node
func getLatestBlock(client *ethclient.Client) (*ethclient.Block, error) {
	latestBlock, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return latestBlock, nil
}

// Check the health of the service
func checkService() bool {
	// TODO: Implement the service health check logic
	return true
}
