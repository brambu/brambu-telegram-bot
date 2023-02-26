package helpers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
)

func ReplyWithText(bot *tgbotapi.BotAPI, update tgbotapi.Update, text string) {
	message := tgbotapi.NewMessage(GetUpdateMessageChatId(update), text)
	message.ParseMode = "markdown"
	message.ReplyToMessageID = GetUpdateMessageMessageId(update)
	_, err := bot.Send(message)
	if err != nil {
		log.Error().Err(err).Str("text", text).Msg("error ReplyWithText")
	}
}
