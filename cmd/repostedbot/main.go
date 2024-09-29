package main

import (
	"io/fs"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/azekeil/repostedbot/internal/config"
	"github.com/azekeil/repostedbot/internal/handler"
	"github.com/azekeil/repostedbot/internal/reposted"
	"github.com/azekeil/repostedbot/internal/self"

	"github.com/bwmarrin/discordgo"
)

func main() {
	v := config.ReadConfig()
	config := config.ParseConfig(v)

	// Setup help
	help := self.MakeHelp("commands", "commands", "Command")

	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = reposted.LoadDB()
	if err != nil {
		if _, ok := err.(*fs.PathError); !ok {
			log.Fatalf("Fatal error loading DB: %v", err)
		}
		log.Printf("Error loading DB: %v", err)
	}
	// Catch up
	err = reposted.CatchUp(session)
	if err != nil {
		log.Printf("Error processing message history: %v", err)
	}

	// Register callback for the messageCreate events.
	session.AddHandler(handler.MakeMessageCreateHandlerFunc(help))

	// Register callback for the guildCreate events.
	session.AddHandler(handler.GuildCreate)

	// Open websocket after registering
	err = session.Open()
	if err != nil {
		panic(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Caught signal!")
	log.Println("Closing connection...")
	err = session.Close()
	if err != nil {
		log.Printf("could not close session gracefully: %s", err)
	}
}
