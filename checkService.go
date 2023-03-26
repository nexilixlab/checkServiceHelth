package main

import (
	"fmt"
	"time"
)

// فراخوانی بخش‌های دیگر برنامه
func main() {
	// خواندن تنظیمات اولیه
	config, err := readConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// برنامه ریزی برای اجرای بخش‌های مختلف هر ۵ دقیقه
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Checking block...")
			blockNumber, err := checkBlock()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			if blockNumber == config.Block {
				fmt.Println("No new block found.")
				continue
			}

			fmt.Println("New block found:", blockNumber)

			if err := writeData(blockNumber); err != nil {
				fmt.Println("Error:", err)
				continue
			}

			if blockNumber == config.Block+1 {
				fmt.Println("Restarting nexilix service...")
				if err := restartService(); err != nil {
					fmt.Println("Error restarting nexilix service:", err)
					continue
				}
			}

			config.Block = blockNumber
			if err := writeConfig(config); err != nil {
				fmt.Println("Error:", err)
				continue
			}
		}
	}
}
