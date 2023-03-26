package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Service struct {
	Name string `json:"name"`
}

type Config struct {
	Jsonrpc string    `json:"jsonrpc"`
	Id      int       `json:"id"`
	Method  string    `json:"method"`
	Params  []Service `json:"params"`
}

const (
	blockCheckInterval = 10 * time.Second
)

func checkBlock(rpcUrl string) bool {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal(err)
		return false
	}

	// Get the latest block number
	latestBlock, err := getLatestBlock(client)
	if err != nil {
		log.Fatal(err)
		return false
	}
	if !checkService() {
		return false
	}
	time.Sleep(blockCheckInterval)
	latestBlockAgain, err := getLatestBlock(client)
	if err != nil {
		log.Fatal(err)
	}

	// If the latest block has changed, log the error and exit the loop
	if latestBlockAgain.NumberU64() == latestBlock.NumberU64() {
		log.Fatalf("Error in service block is %s", latestBlockAgain.NumberU64())
		return false
	}
	return true
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

func readConfig(filename string) (Config, error) {
	var config Config

	configFile, err := os.Open(filename)
	if err != nil {
		return config, err
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}

func writeData(block string) error {
	xmlString := []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><config><lastblock>" + block + "</lastblock></config>")
	err := ioutil.WriteFile("config.xml", xmlString, 0644)
	if err != nil {
		return err
	}

	return nil
}

func restartService(serviceName string, logger *log.Logger) error {
	logger.Printf("Trying to restart service %s\n", serviceName)

	cmd := exec.Command("systemctl", "restart", serviceName)
	err := cmd.Run()

	if err != nil {
		logger.Printf("Failed to restart service %s: %s\n", serviceName, err.Error())
		return err
	}

	logger.Printf("Successfully restarted service %s\n", serviceName)

	return nil
}

func getCurrentBlock(rpcUrl string) {

}

func main() {
	servicesPtr := flag.String("services", "", "comma-separated list of services to restart")
	rpcUrl := flag.String("rpcUrl", "http://127.0.0.1:8545", "")
	flag.Parse()

	services := strings.Split(*servicesPtr, ",")

	if len(services) == 0 {
		//fmt.Println("Error: at least one service name should be provided")
		log.Fatal("Error: at least one service name should be provided")
		return
	}

	logFile, err := os.OpenFile("/var/log/checkNexilix/service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	for {
		time.Sleep(5 * time.Minute)

		if checkBlock(rpcUrl) {
			logger.Println("Service is Helathy")
		} else {
			for _, service := range services {
				restartService(service, logger)
			}
		}
	}
}
