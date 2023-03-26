package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

func checkBlock() (int, error) {
	// ارسال درخواست به آدرس RPC برای گرفتن آخرین بلوک
	resp, err := http.Post("http://127.0.0.1:8545", "application/json", nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// پردازش بلوک دریافتی
	blockNumber, err := strconv.Atoi(string(body))
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}
