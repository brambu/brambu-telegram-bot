package helpers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
)

func SendVoice(bot *tgbotapi.BotAPI, update tgbotapi.Update, data []byte, caption string) {
	if len(data) == 0 {
		log.Error().Msg("SendVoice no data")
		return
	}
	tmpFile := GetTmpFile(data)
	defer CleanupTmpFile(tmpFile)
	message := tgbotapi.NewVoiceUpload(GetUpdateMessageChatId(update), tmpFile.Name())
	message.ReplyToMessageID = GetUpdateMessageMessageId(update)
	message.Caption = caption
	_, err := bot.Send(message)
	if err != nil {
		log.Warn().Err(err).Msg("could not SendVoice")
	}
}
