package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brambu/brambu-telegram-bot/interfaces"
	"net/http"
)

type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func SendMessageToChat(module interfaces.BotModule, chatId int64, messageText string) error {
	config := module.Config()
	baseUrl := "https://api.telegram.org/bot"
	route := "/sendMessage"
	parseMode := "Markdown"
	reqBody := &sendMessageReqBody{
		ChatID: chatId,
		Text:   messageText,
		ParseMode: parseMode,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s%s", baseUrl, config.BotToken, route)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New(" unexpected status " + res.Status)
	}

	return nil
}
