package modules

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"context"
	"fmt"
	"github.com/brambu/brambu-telegram-bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

type Speak struct {
	config config.BotConfiguration
}

func (s *Speak) Name() *string {
	name := "speak"
	return &name
}

func (s *Speak) LoadConfig(conf config.BotConfiguration) {
	s.config = conf
}

func (s Speak) Evaluate(update tgbotapi.Update) bool {
	return strings.HasPrefix(strings.ToLower(update.Message.Text), "/speak")
}

func (s Speak) Execute(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// input like /speak this is how I talk
	// input like /speak .fr Je parle comme ça
	speakText, lang := parseSpeakMsg(update)
	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Error().Err(err)
	}
	resp, err := client.SynthesizeSpeech(ctx, getTtsReq(speakText, lang))
	if resp == nil || err != nil {
		log.Error().Err(err).Msg("SynthesizeSpeech empty response")
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "aroo?")
		_, e := bot.Send(message)
		if e != nil {
			log.Error().Err(err)
		}
		return
	}
	tmpFile, err := os.CreateTemp(os.TempDir(), "brambu-telegram-bot-tts-")
	if err != nil {
		log.Error().Err(err).Msg("cannot create temporary file")
	}
	defer func() {
		if err = os.RemoveAll(tmpFile.Name()); err != nil {
			log.Error().Err(err).Msg("failed to remove temporary file")
		}
	}()
	if _, err = tmpFile.Write(resp.AudioContent); err != nil {
		log.Error().Err(err).Msg("failed to write temporary file")
	}
	message := tgbotapi.NewVoiceUpload(update.Message.Chat.ID, tmpFile.Name())
	message.ReplyToMessageID = update.Message.MessageID
	message.Caption = speakText
	_, err = bot.Send(message)
	if err != nil {
		log.Warn().Err(err).Msg("could not NewAudioShare")
	}
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

func parseSpeakMsg(update tgbotapi.Update) (string, string) {
	speakText := strings.Join(strings.Split(update.Message.Text, " ")[1:], " ")
	possibleLang := strings.Split(speakText, " ")[0]
	possibleMessage := strings.Join(strings.Split(speakText, " ")[1:], " ")
	defaultLang := "en"
	lang := defaultLang
	for _, enabledLang := range enabledLanguages() {
		evalLang := fmt.Sprintf(".%s", enabledLang)
		if possibleLang == evalLang {
			lang = enabledLang
			speakText = possibleMessage
		}
	}
	return speakText, lang
}

func enabledLanguages() []string {
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
