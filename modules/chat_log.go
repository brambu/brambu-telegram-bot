package modules

import (
	"github.com/brambu/brambu-telegram-bot/config"
	"github.com/brambu/brambu-telegram-bot/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type ChatLog struct {
	config config.BotConfiguration
	logger logger.Logger
}

func (c *ChatLog) Name() *string {
	// return nil to avoid logging evaluated lines
	return nil
}

func (c *ChatLog) LoadConfig(conf config.BotConfiguration) {
	c.config = conf
	l := logger.Logger{}
	l.LoadConfig(c.config.LogPath)
	c.logger = l
}

func (c *ChatLog) Evaluate(update tgbotapi.Update) bool {
	return c.config.LogEnabled
}

func (c *ChatLog) Execute(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	c.logger.HandleUpdate(&update)
}
