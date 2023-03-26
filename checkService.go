package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
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

func restartService() error {
	cmd := exec.Command("systemctl", "restart", "nexilix.service")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	for {
		config, err := config.Read()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		blockNumber, err := checkBlock()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if blockNumber != config.Block {
			if err := config.Write(blockNumber); err != nil {
				fmt.Println("Error:", err)
				return
			}

			if err := restartService(); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		time.Sleep(5 * time.Minute)
	}
}
