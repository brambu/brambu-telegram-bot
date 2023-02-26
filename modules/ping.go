package modules

import (
	"github.com/brambu/brambu-telegram-bot/config"
	. "github.com/brambu/brambu-telegram-bot/helpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Ping struct {
	config config.BotConfiguration
}

func (p *Ping) Name() *string {
	name := "ping"
	return &name
}

func (p *Ping) LoadConfig(conf config.BotConfiguration) {
	p.config = conf
}

func (p *Ping) Evaluate(update tgbotapi.Update) bool {
	return CheckPrefix(update, "ping!")
}

func (p *Ping) Execute(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	ReplyWithText(bot, update, "pong!")
}
