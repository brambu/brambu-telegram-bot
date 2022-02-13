package config

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type BotConfiguration struct {
	ConfigPath   string
	BotName      string `yaml:"bot_name"`
	BotToken     string `yaml:"bot_token"`
	DarkskyToken string `yaml:"darksky_token"`
	LogPath      string `yaml:"log_path"`
	LogEnabled   bool   `yaml:"log_enabled"`
	Port         string `yaml:"port"`
	WebhookUrl   string `yaml:"webhook_url"`
	GKeyPath     string `yaml:"gcloud_key_path"`
}

func (botConfig *BotConfiguration) LoadConfiguration(filePath string) *BotConfiguration {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error().Err(err).
			Str("file_path", filePath).
			Msg("error reading config file")
		panic(err)
	}
	botConfig.ConfigPath = filePath
	log.Info().Str("file_path", filePath).Msg("loaded config")
	err = yaml.Unmarshal(yamlFile, &botConfig)
	if err != nil {
		log.Error().Err(err).
			Str("file_path", filePath).
			Msg("error unmarshalling config file")
		panic(err)
	}
	return botConfig
}
