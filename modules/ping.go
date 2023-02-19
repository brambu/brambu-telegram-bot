package modules

import (
	"github.com/brambu/brambu-telegram-bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
	"strings"
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
	return strings.HasPrefix(strings.ToLower(update.Message.Text), "ping!")
}

func (p *Ping) Execute(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	message := tgbotapi.NewMessage(update.Message.Chat.ID, "pong!")
	message.ReplyToMessageID = update.Message.MessageID
	_, err := bot.Send(message)
	if err != nil {
		log.Error().Err(err).Msgf("ping error")
	}
}
