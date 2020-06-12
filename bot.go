package main

import (
	"github.com/brambu/brambu-telegram-bot/bot"
	"github.com/brambu/brambu-telegram-bot/config"
	"github.com/brambu/brambu-telegram-bot/interfaces"
	"github.com/brambu/brambu-telegram-bot/modules"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	NameOfApp = "brambuTelegramBot"
)

func keyboardInterruptHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Interrupted")
		os.Exit(0)
	}()
}

func main() {
	log.Printf("%s starting.", NameOfApp)
	keyboardInterruptHandler()

	log.Println("Loading configuration...")
	var conf config.BotConfiguration
	if len(os.Args) != 2 {
		log.Printf("Add config to command, ex: %s myconfig.yml", os.Args[0])
		os.Exit(1)
	}
	conf.LoadConfiguration(os.Args[1])

	log.Printf("Bot is now %s", conf.BotName)

	modules := []interfaces.BotModule{
		// add modules here
		&botModules.ChatLog{},
		&botModules.Ping{},
	}

	b := bot.WebhookBot{Config: conf, BotModules: modules}
	b.Run()

	log.Printf("Bot %s complete.", conf.BotName)
	log.Printf("%s done, exiting.", NameOfApp)
}
