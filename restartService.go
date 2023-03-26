package main

import (
	"fmt"
	"os/exec"
)

func restartService() error {
	cmd := exec.Command("systemctl", "restart", "nexilix.service")
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println("Nexilix service restarted.")
	return nil
}
