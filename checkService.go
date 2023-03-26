package main

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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

		method := config.Method
		params := config.Params

		url := fmt.Sprintf("%s", config.Jsonrpc)
		payload, err := json.Marshal(map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  method,
			"params":  params,
		})
		if err != nil {
			logger.Printf("Failed to create payload: %s\n", err.Error())
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
		if err != nil {
			logger.Printf("Failed to create HTTP request: %s\n", err.Error())
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Printf("Failed to send HTTP request: %s\n", err.Error())
		}

		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Printf("Failed to read response body: %s\n", err.Error())
			continue
		}

		var respData map[string]interface{}

		if err := json.Unmarshal(respBody, &respData); err != nil {
			logger.Printf("Failed to parse response body: %s\n", err.Error())
			continue
		}

		result, ok := respData["result"].(map[string]interface{})
		if !ok {
			logger.Println("Failed to parse result field in response")
			continue
		}

		block, ok := result["number"].(string)
		if !ok {
			logger.Println("Failed to parse block number from response")
			continue
		}

		lastBlock, err := strconv.Atoi(string(block))
		if err != nil {
			logger.Printf("Failed to parse block number as integer: %s\n", err.Error())
			continue
		}

		lastBlockString := strconv.Itoa(lastBlock)

		xmlFile, err := os.Open("config.xml")
		if err != nil {
			logger.Printf("Failed to open config.xml: %s\n", err.Error())
			continue
		}

		defer xmlFile.Close()

		xmlData, err := ioutil.ReadAll(xmlFile)
		if err != nil {
			logger.Printf("Failed to read config.xml: %s\n", err.Error())
			continue
		}

		var configData map[string]string

		if err := xml.Unmarshal(xmlData, &configData); err != nil {
			logger.Printf("Failed to parse config.xml: %s\n", err.Error())
			continue
		}

		previousBlock, ok := configData["lastblock"]
		if !ok {
			logger.Println("Failed to find last block number in config.xml")
			continue
		}

		if previousBlock == lastBlockString {
			logger.Printf("Server is healthy. Last block is %s and current block is %s\n", previousBlock, lastBlockString)
		} else {
			logger.Printf("Trying to restart services...\n")

			for _, service := range config.Params {
				if err := restartService(service.Name, logger); err != nil {
					logger.Printf("Failed to restart service %s: %s\n", service.Name, err.Error())
				}
			}
		}

		if err := writeData(lastBlockString); err != nil {
			logger.Printf("Failed to write block number to config.xml: %s\n", err.Error())
		}
	}
}
