package bot

import (
	"fmt"
	"github.com/brambu/brambu-telegram-bot/config"
	"github.com/brambu/brambu-telegram-bot/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
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
		log.Fatal(err)
	}
	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(w.Config.WebhookUrl + w.Config.BotToken))
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
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
