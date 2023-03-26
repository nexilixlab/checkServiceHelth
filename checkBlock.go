package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func checkBlock() (int, error) {
	resp, err := http.Get("http://127.0.0.1:8545")
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	blockNumberHex := result["result"].(string)
	blockNumber, err := strconv.ParseInt(blockNumberHex[2:], 16, 64)
	if err != nil {
		return 0, err
	}

	return int(blockNumber), nil
}
