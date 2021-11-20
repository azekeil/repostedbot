package main

import (
	"go/doc"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/alex-broad/grec/internal/handler"
	"github.com/alex-broad/grec/internal/self"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

var rootPath string

func getRootPath() string {
	if rootPath == "" {
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		rootPath = filepath.Join(basepath, "../../")
	}
	return rootPath
}

func readConfig() *viper.Viper {
	v := viper.New()
	v.AddConfigPath(getRootPath())
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return v
}

type config struct {
	Token string `yaml:"token"`
}

func parseConfig(v *viper.Viper) *config {
	c := &config{}
	err := v.Unmarshal(c)
	if err != nil {
		panic(err)
	}
	return c
}

func setupHelp(path, pkg, typ string) map[string]*doc.Func {
	relpath := filepath.Join(getRootPath(), path)

	d, err := self.ParseDir(relpath)
	if err != nil {
		panic(err)
	}
	return self.GetDocMethods(self.GetDocPackage(d, pkg), typ)
}

func main() {
	v := readConfig()
	config := parseConfig(v)

	// Setup help
	help := setupHelp("internal/commands", "commands", "Command")

	bot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		panic(err)
	}
	defer bot.Close()

	bot.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Register messageCreate as a callback for the messageCreate events.
	bot.AddHandler(handler.MakeMessageCreateHandlerFunc(help))

	// Register guildCreate as a callback for the guildCreate events.
	bot.AddHandler(guildCreate)

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

// This function will be called (due to AddHandler above) every time a new
// guild is joined.
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, err := s.ChannelMessageSend(channel.ID, "Hi! Use \"!grec\" to get started")
			if err != nil {
				log.Printf("could not send guild creation message: %s", err)
			}
			log.Printf("Guild %s greeted", event.Name)
			return
		}
	}
}
