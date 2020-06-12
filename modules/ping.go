package botModules

import (
	"brambu-telegram-bot/config"
	"brambu-telegram-bot/helpers"
	"log"
	"strings"
)

type Ping struct {
	config config.BotConfiguration
}

func (p *Ping) LoadConfig(conf config.BotConfiguration) {
	p.config = conf
}

func (p *Ping) Config() config.BotConfiguration {
	return p.config
}

func (p Ping) Evaluate(chatId int64, messageText string, raw string) bool {
	if strings.Contains(strings.ToLower(messageText), "ping!") {
		log.Println("Ping detected.")
		return true
	}
	return false
}

func (p Ping) Execute(chatId int64, messageText string, raw string) {
	log.Println("Sending pong.")
	err := helpers.SendMessageToChat(&p, chatId, "pong!")
	if err != nil {
		log.Printf("Warning: Ping error %s", err)
	}
}