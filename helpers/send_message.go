package helpers

import (
	"brambu-telegram-bot/interfaces"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func SendMessageToChat(module interfaces.BotModule, chatId int64, messageText string) error {
	config := module.Config()
	baseUrl := "https://api.telegram.org/bot"
	route := "/sendMessage"
	reqBody := &sendMessageReqBody{
		ChatID: chatId,
		Text:   messageText,
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
