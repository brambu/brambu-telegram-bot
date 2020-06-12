package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brambu/brambu-telegram-bot/config"
	"github.com/brambu/brambu-telegram-bot/interfaces"
	"log"
	"net/http"
)

type WebhookBot struct {
	Config     config.BotConfiguration
	BotModules []interfaces.BotModule
}

type webhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID       int64  `json:"id"`
			Username string `json:"username"`
		} `json:"chat"`
	} `json:"message"`
}

func (w *WebhookBot) bootstrapModules() {
	for _, module := range w.BotModules {
		module.LoadConfig(w.Config)
	}
}

func (w WebhookBot) Handler(res http.ResponseWriter, req *http.Request) {
	body := &webhookReqBody{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	reqBody := buf.String()
	log.Println(reqBody)
	if err := json.NewDecoder(buf).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}
	for _, module := range w.BotModules {
		if module.Evaluate(body.Message.Chat.ID, body.Message.Text, reqBody) == true {
			module.Execute(body.Message.Chat.ID, body.Message.Text, reqBody)
		}
		return
	}
}

func (w WebhookBot) Run() {
	w.bootstrapModules()
	port := fmt.Sprintf(":%s", w.Config.Port)
	http.ListenAndServe(port, http.HandlerFunc(w.Handler))
}
