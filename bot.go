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
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Error().Msg("Interrupted")
		os.Exit(0)
	}()
}

func main() {
	log.Info().Str("name_of_app", NameOfApp).Msg("starting")
	keyboardInterruptHandler()

	log.Info().Msg("loading configuration")
	var conf config.BotConfiguration
	if len(os.Args) != 2 {
		log.Error().Msgf("add config to command, ex: %s myconfig.yml", os.Args[0])
		os.Exit(1)
	}
	conf.LoadConfiguration(os.Args[1])
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", conf.GKeyPath)
	if err != nil {
		log.Error().Err(err).Msg("bot error setting google application credentials")
		panic(err)
	}

	log.Info().Str("bot_name", conf.BotName).Msg("bot name set")

	botModules := []interfaces.BotModule{
		// add modules here
		&modules.ChatLog{},
		&modules.Ping{},
		&modules.Weather{},
		&modules.Speak{},
	}

	b := bot.WebhookBot{Config: conf, BotModules: botModules}
	err = b.Run()
	if err != nil {
		log.Error().Err(err).Msg("bot error")
		panic(err)
	}
	log.Info().Str("bot_name", conf.BotName).Msg("complete")
	log.Info().Str("name_of_app", NameOfApp).Msg("done, exiting")
}
