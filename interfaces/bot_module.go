package interfaces

import "brambu-telegram-bot/config"

type BotModule interface {
	Evaluate(chatId int64, text string, raw string) bool
	Execute(chatId int64, text string, raw string)
	LoadConfig(config config.BotConfiguration)
	Config() config.BotConfiguration
}