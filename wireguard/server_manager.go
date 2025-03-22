package wireguard

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	wgConfigPath    = "/configs/wg_confs/wg0.conf"
	wgInterface     = "wg0"
	wgRestartScript = "/shell/restart_wireguard.sh"
)

func AddClientToServer(telegramID int64, publicKey string, clientIP string) error {
	clientConfig := fmt.Sprintf("\n[Peer]\n# TelegramID: %d\nPublicKey = %s\nAllowedIPs = %s\n", telegramID, publicKey, clientIP)

	f, err := os.OpenFile(wgConfigPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open WireGuard config: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(clientConfig); err != nil {
		return fmt.Errorf("failed to write to WireGuard config: %v", err)
	}

	if err := restartWireGuard(); err != nil {
		return fmt.Errorf("failed to restart WireGuard: %v", err)
	}

	return nil
}



func RemoveClientServer(telegramID int64) error {
	configData, err := os.ReadFile(wgConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read WireGuard config: %v", err)
	}

	configString := string(configData)
	peerConfig := fmt.Sprintf("[Peer]\n# TelegramID: %d", telegramID)

	for i := 0; i < 3; i++ {
		newConfig := removePeerBlock(configString, peerConfig)
		if newConfig == configString {
			break
		}
		configString = newConfig
	}

	if err := os.WriteFile(wgConfigPath, []byte(configString), 0600); err != nil {
		return fmt.Errorf("failed to update WireGuard config: %v", err)
	}

	if err := restartWireGuard(); err != nil {
		return fmt.Errorf("failed to restart WireGuard after removal: %v", err)
	}

	log.Printf("Client with TelegramID %d removed from WireGuard", telegramID)
	return nil
}


func removePeerBlock(config, peer string) string {
	idx := strings.Index(config, peer)
	if idx == -1 {
		return config
	}

	startIdx := idx
	endIdx := strings.Index(config[startIdx+len(peer):], "[Peer]")

	if endIdx != -1 {
		endIdx += startIdx + len(peer)
	} else {
		endIdx = len(config)
	}

	if startIdx >= len(config) || endIdx > len(config) {
		return config
	}

	updatedConfig := config[:startIdx] + config[endIdx:]
	return strings.TrimSpace(updatedConfig) + "\n"
}

func restartWireGuard() error {
	var stderr bytes.Buffer

	cmd := exec.Command("/bin/sh", wgRestartScript)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Printf("WireGuard restart failed: %v\nStderr: %s", err, stderr.String())
		return fmt.Errorf("failed to execute restart script: %v", err)
	}
	log.Println("WireGuard restarted successfully.")
	return nil
}