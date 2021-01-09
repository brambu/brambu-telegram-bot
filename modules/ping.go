package modules

import (
	"github.com/brambu/brambu-telegram-bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

type Ping struct {
	config config.BotConfiguration
}

func (p *Ping) LoadConfig(conf config.BotConfiguration) {
	p.config = conf
}

func (p Ping) Evaluate(update tgbotapi.Update) bool {
	if strings.Contains(strings.ToLower(update.Message.Text), "ping!") {
		log.Printf("Ping detected from [%d,%s].",
			update.Message.From.ID,
			update.Message.From.UserName,
		)
		return true
	}
	return false
}

func (p Ping) Execute(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Println("Sending pong.")
	message := tgbotapi.NewMessage(update.Message.Chat.ID, "pong!")
	message.ReplyToMessageID = update.Message.MessageID
	_, err := bot.Send(message)
	if err != nil {
		log.Printf("Warning: Ping error %s", err)
	}
}
