package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/nexilixlab/checkServiceHelth/config"
)

func checkBlock() (int, error) {
	resp, err := http.Get("http://127.0.0.1:8545")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	blockNumberHex := result["result"].(string)
	blockNumber, err := strconv.ParseInt(blockNumberHex[2:], 16, 64)
	if err != nil {
		return 0, err
	}

	return int(blockNumber), nil
}

func restartService(service string) error {
	cmd := exec.Command("systemctl", "restart", service)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	servicesPtr := flag.String("services", "", "comma-separated list of services to restart")
	flag.Parse()

	services := strings.Split(*servicesPtr, ",")

	if len(services) == 0 {
		fmt.Println("Error: at least one service name should be provided")
		return
	}
	for {
		conf, err := config.Read()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		blockNumber, err := checkBlock()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if blockNumber != conf.Block {
			if err := config.Write(blockNumber); err != nil {
				fmt.Println("Error: ", err)
				return
			}

			for _, service := range services {
				if err := restartService(service); err != nil {
					fmt.Println("Error: ", err)
					return
				}
			}
		}

		time.Sleep(5 * time.Minute)
	}
}
