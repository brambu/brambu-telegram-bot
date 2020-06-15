package modules

import (
	"encoding/json"
	"fmt"
	"github.com/brambu/brambu-telegram-bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

type ChatLog struct {
	config     config.BotConfiguration
	fileHandle *os.File
	logger     *log.Logger
}

func (c *ChatLog) LoadConfig(conf config.BotConfiguration) {
	c.config = conf
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

func (c ChatLog) Evaluate(update tgbotapi.Update) bool {
	if c.config.LogEnabled == true {
		return true
	}
	return false
}

func (c ChatLog) Execute(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	message, err := json.Marshal(update.Message)
	if err != nil {
		log.Printf("ChatLog marshal error %s", err)
		return
	}
	stringMessage := fmt.Sprintf("%s", message)
	err = c.LogLineToFile(stringMessage)
	if err != nil {
		log.Printf("ChatLog log error %s", err)
	}
}
