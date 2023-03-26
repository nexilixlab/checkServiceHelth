package config

import (
	"encoding/xml"
	"io/ioutil"
)

// تعریف ساختار داده‌ای برای فایل تنظیمات
type Config struct {
	Block int `xml:"block"`
}

// تابع خواندن تنظیمات از فایل XML
func Read() (*Config, error) {
	file, err := ioutil.ReadFile("config.xml")
	if err != nil {
		return nil, err
	}

	var config Config
	if err := xml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// تابع نوشتن آخرین بلوک در فایل XML
func Write(blockNumber int) error {
	config := Config{Block: blockNumber}
	data, err := xml.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("config.xml", []byte(xml.Header+string(data)), 0644); err != nil {
		return err
	}

	return nil
}
