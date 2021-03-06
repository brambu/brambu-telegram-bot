package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
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
	GKeyPath	 string `yaml:"gcloud_key_path"`
}

func (botConfig *BotConfiguration) LoadConfiguration(filePath string) *BotConfiguration {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading config file %s #%v", filePath, err)
		os.Exit(2)
	}
	botConfig.ConfigPath = filePath
	log.Printf("Loaded :%s", filePath)
	err = yaml.Unmarshal(yamlFile, &botConfig)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return botConfig
}
