package interfaces

import (
	"github.com/brambu/brambu-telegram-bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotModule interface {
	Evaluate(update tgbotapi.Update) bool
	Execute(bot *tgbotapi.BotAPI, update tgbotapi.Update)
	LoadConfig(config config.BotConfiguration)
	Config() config.BotConfiguration
}
