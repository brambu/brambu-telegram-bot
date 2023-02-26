package logger

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	zlog "github.com/rs/zerolog/log"
	"log"
	"os"
)

type Logger struct {
	path   string
	logger *log.Logger
}

func (l *Logger) LoadConfig(path string) {
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		zlog.Error().Err(err).Str("path", path).Msg("logger error opening path")
	}
	l.logger = log.New(f, "", log.LstdFlags)
	l.path = path
}

func (l *Logger) HandleUpdate(update *tgbotapi.Update) {
	if update != nil {
		if update.Message != nil {
			l.LogMessage(update.Message)
		}
	}
}

func (l *Logger) LogMessage(message *tgbotapi.Message) {
	bytes, err := json.Marshal(message)
	if err != nil {
		zlog.Error().Err(err).Msg("could not marshal log message")
	}
	l.logger.Println(string(bytes))
}
