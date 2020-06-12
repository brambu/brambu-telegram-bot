package modules

import (
	"fmt"
	"github.com/brambu/brambu-telegram-bot/config"
	"github.com/brambu/brambu-telegram-bot/helpers"
	"github.com/evalphobia/google-tts-go/googletts"
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

func (s Speak) Evaluate(chatId int64, messageText string, raw string) bool {
	if strings.HasPrefix(strings.ToLower(messageText), "/speak") {
		log.Printf("Speak command: %s", messageText)
		return true
	}
	return false
}

func (s Speak) Execute(chatId int64, messageText string, raw string) {
	log.Println("Sending speak response.")
	// input like /speak this is how I talk
	// input like /speak .fr Je parle comme ça
	speakText := strings.Join(strings.Split(messageText, " ")[1:], " ")
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
	err = helpers.SendAudioToChatByURL(&s, chatId, url, speakText)
	if err != nil {
		log.Printf("Warning: Speak error %s", err)
	}
}
