package main

import (
	"encoding/xml"
	"io/ioutil"
)

func writeData(blockNumber int) error {
	config := &Config{Block: blockNumber}

	// تبدیل تنظیمات به داده‌های XML
	data, err := xml.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// ذخیره تنظیمات در فایل
	if err := ioutil.WriteFile("config.xml", data, 0644); err != nil {
		return err
	}

	return nil
}
