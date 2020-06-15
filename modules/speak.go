package modules

import (
	"fmt"
	"github.com/brambu/brambu-telegram-bot/config"
	"github.com/evalphobia/google-tts-go/googletts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

type Speak struct {
	config config.BotConfiguration
}

func (s *Speak) LoadConfig(conf config.BotConfiguration) {
	s.config = conf
}

func (s *Speak) Config() config.BotConfiguration {
	return s.config
}

func (s Speak) EnabledLanguages() []string {
	return []string{
		"ar",
		"bn",
		"cmn",
		"da",
		"de",
		"el",
		"en",
		"es",
		"fi",
		"fil",
		"fr",
		"gu",
		"hi",
		"hu",
		"id",
		"it",
		"ja",
		"kn",
		"ml",
		"nb",
		"nl",
		"pl",
		"pt",
		"ru",
		"sk",
		"ta",
		"te",
		"th",
		"tr",
		"uk",
		"vi",
	}
}

func (s Speak) Evaluate(update tgbotapi.Update) bool {
	if strings.HasPrefix(strings.ToLower(update.Message.Text), "/speak") {
		log.Printf("Speak command: %s", update.Message.Text)
		return true
	}
	return false
}

func (s Speak) Execute(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Println("Sending speak response.")
	// input like /speak this is how I talk
	// input like /speak .fr Je parle comme ça
	speakText := strings.Join(strings.Split(update.Message.Text, " ")[1:], " ")
	possibleLang := strings.Split(speakText, " ")[0]
	possibleMessage := strings.Join(strings.Split(speakText, " ")[1:], " ")
	defaultLang := "en"
	lang := defaultLang
	for _, enabledLang := range s.EnabledLanguages() {
		// require user to specify a dot before the language
		evalLang := fmt.Sprintf(".%s", enabledLang)
		if possibleLang == evalLang {
			lang = enabledLang
			speakText = possibleMessage
		}
	}
	url, err := googletts.GetTTSURL(speakText, lang)
	if err != nil {
		log.Printf("Warning: Speak GetTTS error %s", err)
	}

	message := tgbotapi.NewAudioShare(update.Message.Chat.ID, url)
	_, err = bot.Send(message)
	if err != nil {
		log.Printf("Warning: could not NewAudioShare %s", err)
	}
}
