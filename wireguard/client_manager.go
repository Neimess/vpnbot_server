package wireguard

import (
	"fmt"
	"os"
	"strings"

	"github.com/Neimess/vpnbot_server/config"
	"gorm.io/gorm"
)

const (
	wgConfigDir = "/configs/clients/"
)

func GenerateClientConfig(db *gorm.DB, clientIP string, telegramID int64, privateKey string) (string, string, error) {
	config := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s
DNS = 1.1.1.1

[Peer]
PublicKey = %s
Endpoint = %s:%s
AllowedIPs = 0.0.0.0/0
PersistentKeepalive = 25
`, privateKey,
		clientIP,
		config.GlobalConfig.SERVER_PUBLIC_KEY,
		config.GlobalConfig.SERVER_IP,
		config.GlobalConfig.WG_PORT)

	lastOctet := strings.Split(clientIP, ".")[3]
	lastOctet = strings.Split(lastOctet, "/")[0]
	filePath := fmt.Sprintf("%sclient_%s.conf", wgConfigDir, lastOctet)
	err := os.WriteFile(filePath, []byte(config), 0600)
	if err != nil {
		return "", "", fmt.Errorf("failed to write config file: %v", err)
	}

	return config, filePath, nil
}

func RemoveClientData(config_path string) error {
	err := os.Remove(config_path)
	if err != nil {
		return err
	}
	return nil
}
