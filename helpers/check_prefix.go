package helpers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

func CheckPrefix(update tgbotapi.Update, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(GetUpdateMessageText(update)), prefix)
}
