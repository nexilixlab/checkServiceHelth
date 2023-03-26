package readConfig

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

// تعریف ساختار داده‌ای برای فایل تنظیمات
type Config struct {
	Block int `xml:"block"`
}

// تابع خواندن تنظیمات از فایل XML
func Read() (*Config, error) {
	file, err := os.Open("config.xml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := xml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
