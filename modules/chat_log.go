package botModules

import (
	"brambu-telegram-bot/config"
	"log"
	"os"
)

type ChatLog struct {
	config config.BotConfiguration
	fileHandle *os.File
	logger *log.Logger
}

func (c *ChatLog) LoadConfig(conf config.BotConfiguration) {
	c.config = conf
}

func (c *ChatLog) Config() config.BotConfiguration {
	return c.config
}

func (c *ChatLog) LogLineToFile(line string) error {
	if c.fileHandle == nil {
		f, err := os.OpenFile(c.config.LogPath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		c.fileHandle = f
		defer c.fileHandle.Close()
	}
	if c.logger == nil {
		logger := log.New(c.fileHandle, "", log.LstdFlags)
		c.logger = logger
	}
	c.logger.Println(line)
	return nil
}

func (c ChatLog) Evaluate(chatId int64, messageText string, raw string) bool {
	if c.config.LogEnabled == true {
		return true
	}
	return false
}

func (c ChatLog) Execute(chatId int64, messageText string, raw string) {
	c.LogLineToFile(raw)
}