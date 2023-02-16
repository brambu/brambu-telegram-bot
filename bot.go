package main

import (
	"github.com/brambu/brambu-telegram-bot/bot"
	"github.com/brambu/brambu-telegram-bot/config"
	"github.com/brambu/brambu-telegram-bot/interfaces"
	"github.com/brambu/brambu-telegram-bot/modules"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

const (
	NameOfApp = "BrambuTelegramBot"
)

func keyboardInterruptHandler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		log.Error().Msg("Interrupted")
		os.Exit(0)
	}()
}

func loadConfig() config.BotConfiguration {
	log.Info().Msg("loading configuration")
	var conf config.BotConfiguration
	if len(os.Args) != 2 {
		log.Error().Msgf("add config to command, ex: %s myconfig.yml", os.Args[0])
		os.Exit(1)
	}
	conf.LoadConfiguration(os.Args[1])
	log.Info().Str("bot_name", conf.BotName).Msg("bot name set")
	return conf
}

func setGoogCreds(conf config.BotConfiguration) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", conf.GKeyPath)
	if err != nil {
		log.Error().Err(err).Msg("bot error setting google application credentials")
		panic(err)
	}
}

func main() {
	log.Info().Str("name_of_app", NameOfApp).Msg("starting")
	keyboardInterruptHandler()
	conf := loadConfig()
	setGoogCreds(conf)
	botModules := []interfaces.BotModule{
		// add modules here
		&modules.ChatLog{},
		&modules.Ping{},
		&modules.Weather{},
		&modules.Speak{},
	}
	b := bot.WebhookBot{Config: conf, BotModules: botModules}
	err := b.Run()
	if err != nil {
		log.Error().Err(err).Msg("bot error")
		panic(err)
	}
	log.Info().Str("bot_name", conf.BotName).Msg("complete")
	log.Info().Str("name_of_app", NameOfApp).Msg("done, exiting")
}
