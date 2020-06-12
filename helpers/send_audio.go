package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brambu/brambu-telegram-bot/interfaces"
	"net/http"
)

type sendAudioReqBody struct {
	ChatID int64  `json:"chat_id"`
	Audio   string `json:"audio"`
	Caption string `json:"caption"`
	ParseMode string `json:"parse_mode"`
}

func SendAudioToChatByURL(module interfaces.BotModule, chatId int64, url string, caption string) error {
	config := module.Config()
	baseUrl := "https://api.telegram.org/bot"
	route := "/sendAudio"
	parseMode := "Markdown"
	reqBody := &sendAudioReqBody{
		ChatID: chatId,
		Caption: caption,
		Audio: url,
		ParseMode: parseMode,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	u := fmt.Sprintf("%s%s%s", baseUrl, config.BotToken, route)
	res, err := http.Post(u, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New(" unexpected status " + res.Status)
	}

	return nil
}
