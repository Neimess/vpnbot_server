package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Neimess/vpnbot_server/config"
)

type TelegramSendMessageRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func SendTelegramMessage(telegramID int64, message string) error {
	botToken := config.GlobalConfig.TELEGRAM_BOT_TOKEN
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	payload := TelegramSendMessageRequest{
		ChatID: telegramID,
		Text:   message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}
	return nil
}
