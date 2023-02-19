package modules

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"context"
	"fmt"
	"github.com/brambu/brambu-telegram-bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
	"sort"
	"strings"
)

type Speak struct {
	config           config.BotConfiguration
	enabledLanguages []string
	ctx              context.Context
	client           *texttospeech.Client
}

func (s *Speak) Name() *string {
	name := "speak"
	return &name
}

func (s *Speak) LoadConfig(conf config.BotConfiguration) {
	s.config = conf
	s.ctx = context.Background()
	client, err := texttospeech.NewClient(s.ctx)
	if err != nil {
		log.Error().Err(err)
		panic(err)
	}
	s.client = client
	s.getEnabledLanguages()
}

func (s *Speak) Evaluate(update tgbotapi.Update) bool {
	return strings.HasPrefix(strings.ToLower(update.Message.Text), "/speak")
}

func (s *Speak) Execute(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// input like /speak this is how I talk
	// input like /speak .fr Je parle comme ça
	speakText, lang := parseSpeakMsg(update, s.enabledLanguages)
	if lang == "help" {
		sendMessage(bot, update, s.getHelpMsg())
		return
	}
	if speakText == "" {
		sendMessage(bot, update, "aroo?")
		return
	}
	resp, err := s.client.SynthesizeSpeech(s.ctx, getTtsReq(speakText, lang))
	if resp == nil || err != nil {
		log.Error().Err(err).Msg("SynthesizeSpeech empty response")
		sendMessage(bot, update, "aroo?")
		return
	}
	file := tgbotapi.FileBytes{
		Name:  speakText,
		Bytes: resp.GetAudioContent(),
	}
	message := tgbotapi.NewVoiceUpload(update.Message.Chat.ID, file)
	message.ReplyToMessageID = update.Message.MessageID
	message.Caption = speakText
	_, err = bot.Send(message)
	if err != nil {
		log.Warn().Err(err).Msg("could not NewAudioShare")
	}
}

func (s *Speak) getHelpMsg() string {
	return fmt.Sprintf(
		"examples:\n\n%s\n%s\n\n[pick a lanugage code from here](%s)",
		"`/speak this is how I talk`",
		"`/speak .fr Je parle comme ça`",
		"https://cloud.google.com/text-to-speech/docs/voices",
	)
}

func (s *Speak) getEnabledLanguages() {
	req := new(texttospeechpb.ListVoicesRequest)
	resp, err := s.client.ListVoices(s.ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("error listing voices`")
	}
	tracker := make(map[string]bool)
	for _, voice := range resp.GetVoices() {
		for _, languageCode := range voice.LanguageCodes {
			shortName := strings.Split(languageCode, "-")[0]
			tracker[shortName] = true
			tracker[languageCode] = true
		}
	}
	var languageCodes []string
	for k := range tracker {
		languageCodes = append(languageCodes, k)
	}
	sort.Strings(languageCodes)
	s.enabledLanguages = languageCodes
}

func sendMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, message string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "markdown"
	_, err := bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Str("message", message).Msg("could not send speak msg")
	}
}

func parseSpeakMsg(update tgbotapi.Update, enabledLanguages []string) (string, string) {
	speakText := strings.Join(strings.Split(update.Message.Text, " ")[1:], " ")
	possibleLang := strings.Split(speakText, " ")[0]
	if possibleLang == ".help" {
		return "", "help"
	}
	possibleMessage := strings.Join(strings.Split(speakText, " ")[1:], " ")
	defaultLang := "en"
	lang := defaultLang
	for _, enabledLang := range enabledLanguages {
		evalLang := fmt.Sprintf(".%s", enabledLang)
		if possibleLang == evalLang {
			lang = enabledLang
			speakText = possibleMessage
		}
	}
	return speakText, lang
}

func getTtsReq(speakText string, lang string) *texttospeechpb.SynthesizeSpeechRequest {
	return &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: speakText},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: lang,
			SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_OGG_OPUS,
		},
	}
}
