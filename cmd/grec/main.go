package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alex-broad/grec/internal/config"
	"github.com/alex-broad/grec/internal/handler"
	"github.com/alex-broad/grec/internal/self"

	"github.com/bwmarrin/discordgo"
)

func main() {
	v := config.ReadConfig()
	config := config.ParseConfig(v)

	// Setup help
	help := self.MakeHelp("commands", "commands", "Command")

	bot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		panic(err)
	}
	defer bot.Close()

	bot.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Register messageCreate as a callback for the messageCreate events.
	bot.AddHandler(handler.MakeMessageCreateHandlerFunc(help))

	// Register guildCreate as a callback for the guildCreate events.
	bot.AddHandler(handler.GuildCreate)

	// Open websocket after registering
	err = bot.Open()
	if err != nil {
		panic(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
