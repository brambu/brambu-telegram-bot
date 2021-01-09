package modules

import (
	"context"
	"fmt"
	"github.com/brambu/brambu-telegram-bot/config"
	"io/ioutil"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

type Speak struct {
	config config.BotConfiguration
}

func (s *Speak) LoadConfig(conf config.BotConfiguration) {
	s.config = conf
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

	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: speakText},
		},
		// Build the voice request, select the language code
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: lang,
			SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_OGG_OPUS,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "brambu-telegram-bot-tts-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	// Remember to clean up the file afterwards
	defer os.Remove(tmpFile.Name())

	fmt.Println("Created File: " + tmpFile.Name())

	if _, err = tmpFile.Write(resp.AudioContent); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	message := tgbotapi.NewVoiceUpload(update.Message.Chat.ID, tmpFile.Name())
	message.Caption = speakText
	_, err = bot.Send(message)
	if err != nil {
		log.Printf("Warning: could not NewAudioShare %s", err)
	}

	// Close the file
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}
}
