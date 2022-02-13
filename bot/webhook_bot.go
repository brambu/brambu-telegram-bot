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

func (w *WebhookBot) bootstrapModules() {
	for _, module := range w.BotModules {
		go module.LoadConfig(w.Config)
	}
}

func (w WebhookBot) RunModule(bot *tgbotapi.BotAPI, module interfaces.BotModule, update tgbotapi.Update) {
	if module.Evaluate(update) == true {
		module.Execute(bot, update)
	}
}

func (w WebhookBot) Run() error {
	w.bootstrapModules()

	bot, err := tgbotapi.NewBotAPI(w.Config.BotToken)
	if err != nil {
		log.Error().Err(err)
		panic(err)
	}
	// bot.Debug = true
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
	port := fmt.Sprintf(":%s", w.Config.Port)
	go http.ListenAndServe(port, nil)
	for update := range updates {
		for _, module := range w.BotModules {
			go w.RunModule(bot, module, update)
		}
	}
	return nil
}
