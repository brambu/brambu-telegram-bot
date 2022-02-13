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

func (p *Ping) LoadConfig(conf config.BotConfiguration) {
	p.config = conf
}

func (p Ping) Evaluate(update tgbotapi.Update) bool {
	if strings.HasPrefix(strings.ToLower(update.Message.Text), "ping!") {
		log.Info().
			Int("from_id", update.Message.From.ID).
			Str("from_user_name", update.Message.From.UserName).
			Msg("ping detected")
		return true
	}
	return false
}

func (p Ping) Execute(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Info().
		Str("chat_title", update.Message.Chat.Title).
		Int64("chat_id", update.Message.Chat.ID).
		Msg("Sending pong.")
	message := tgbotapi.NewMessage(update.Message.Chat.ID, "pong!")
	message.ReplyToMessageID = update.Message.MessageID
	_, err := bot.Send(message)
	if err != nil {
		log.Error().Err(err).Msgf("ping error")
	}
}
