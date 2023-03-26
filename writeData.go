package main

import (
	"encoding/xml"
	"io/ioutil"
)

const filename = "config.xml"

func writeData(blockNumber int) error {
	config := Config{Block: blockNumber}
	data, err := xml.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, []byte(xml.Header+string(data)), 0644)
	if err != nil {
		return err
	}

	return nil
}
