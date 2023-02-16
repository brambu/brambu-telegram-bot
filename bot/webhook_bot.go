package bot

import (
	"fmt"
	"github.com/brambu/brambu-telegram-bot/config"
	"github.com/brambu/brambu-telegram-bot/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
	"net/http"
)

type WebhookBot struct {
	Config     config.BotConfiguration
	BotModules []interfaces.BotModule
}

func (w WebhookBot) Run() error {
	w.bootstrapModules()
	bot, err := tgbotapi.NewBotAPI(w.Config.BotToken)
	if err != nil {
		log.Error().Err(err)
		panic(err)
	}
	log.Info().Str("bot_self_user_name", bot.Self.UserName).Msg("authorized account")
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(w.Config.WebhookUrl + w.Config.BotToken))
	if err != nil {
		log.Error().Err(err)
		panic(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Error().Err(err)
		panic(err)
	}
	if info.LastErrorDate != 0 {
		log.Error().
			Str("last_error_message", info.LastErrorMessage).
			Msg("telegram callback failed")
	}
	updates := bot.ListenForWebhook("/" + w.Config.BotToken)
	go func() {
		port := fmt.Sprintf(":%s", w.Config.Port)
		err = http.ListenAndServe(port, nil)
		if err != nil {
			log.Error().Err(err)
		}
	}()
	for update := range updates {
		for _, module := range w.BotModules {
			go runModule(bot, module, update)
		}
	}
	return nil
}

func (w *WebhookBot) bootstrapModules() {
	for _, module := range w.BotModules {
		go module.LoadConfig(w.Config)
	}
}

func runModule(bot *tgbotapi.BotAPI, module interfaces.BotModule, update tgbotapi.Update) {
	name := module.Name()
	if module.Evaluate(update) {
		if name != nil {
			log.Info().
				Int("from_id", update.Message.From.ID).
				Str("from_user_name", update.Message.From.UserName).
				Str("text", update.Message.Text).
				Str("chat_title", update.Message.Chat.Title).
				Str("module_name", *name).
				Msg("module exec")
		}
		module.Execute(bot, update)
	}
}
