package helpers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func GetUpdateMessageText(update tgbotapi.Update) string {
	var ret string
	if update.Message != nil {
		ret = update.Message.Text
	}
	return ret
}

func GetUpdateMessageChatId(update tgbotapi.Update) int64 {
	var ret int64
	if update.Message != nil {
		if update.Message.Chat != nil {
			ret = update.Message.Chat.ID
		}
	}
	return ret
}

func GetUpdateMessageMessageId(update tgbotapi.Update) int {
	var ret int
	if update.Message != nil {
		ret = update.Message.MessageID
	}
	return ret
}
