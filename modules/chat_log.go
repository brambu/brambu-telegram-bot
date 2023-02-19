package modules

import (
	"encoding/json"
	"github.com/brambu/brambu-telegram-bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
	golog "log"
	"os"
)

type ChatLog struct {
	config config.BotConfiguration
	logger *golog.Logger
}

func (c *ChatLog) Name() *string {
	// return nil to avoid logging evaluated lines
	return nil
}

func (c *ChatLog) LoadConfig(conf config.BotConfiguration) {
	c.config = conf
}

func (c *ChatLog) Evaluate(update tgbotapi.Update) bool {
	return c.config.LogEnabled
}

func (c *ChatLog) Execute(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	message, err := json.Marshal(update.Message)
	if err != nil {
		log.Error().Err(err).Msg("ChatLog marshal error")
		return
	}
	if err = c.logLineToFile(string(message)); err != nil {
		log.Error().Err(err).Msg("ChatLog log error")
	}
}

func (c *ChatLog) logLineToFile(line string) error {
	if c.logger == nil {
		f, err := os.OpenFile(c.config.LogPath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Error().Err(err)
		}
		logger := golog.New(f, "", golog.LstdFlags)
		c.logger = logger
	}
	c.logger.Println(line)
	return nil
}
