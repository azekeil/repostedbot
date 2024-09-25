package handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/azekeil/repostedbot/internal/bot"
	"github.com/azekeil/repostedbot/internal/commands"
	"github.com/azekeil/repostedbot/internal/self"
	"github.com/bwmarrin/discordgo"
)

func MakeMessageCreateHandlerFunc(help self.DocFuncs) func(*discordgo.Session, *discordgo.MessageCreate) {
	c := new(commands.Command)
	cmdNotFoundFmtStr := "Could not find command `%s`. Try `!rp list`"
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {

		// Ignore all messages created by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Ignore any messages that don't start with "!rp"
		if !strings.HasPrefix(m.Content, "!rp") {
			return
		}

		all := strings.Fields(m.Content)
		if len(all) < 2 {
			c.Help(s, m, help)
			return
		}

		// OK, we know there must be a command..
		// See if there's a function with the same name, if so, call it
		msg := ""
		msgErr := false
		// Special case for the 'help <x>' function
		if all[1] == "help" && len(all) > 2 {
			if cmdhelp := help.Capitalise(all[2]); cmdhelp != "" {
				msg = help.CommandHelp(cmdhelp)
			} else {
				msg = fmt.Sprintf(cmdNotFoundFmtStr, all[2])
				msgErr = true
			}
		} else if command := help.Capitalise(all[1]); command != "" {
			self.CallMethod(c, command, []interface{}{s, m, help})
		} else {
			msg = fmt.Sprintf(cmdNotFoundFmtStr, all[1])
			msgErr = true
		}

		// Send message if defined
		if msg != "" {
			fn := bot.NewEmbed
			if msgErr {
				fn = bot.NewErrorEmbed
			}
			bot.SendEmbed(s, m.ChannelID, fn(msg))
		}
	}
}

// This function will be called (due to AddHandler in main) every time a new
// guild is joined.
func GuildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, err := s.ChannelMessageSend(channel.ID, "Hi! Use \"!rp\" to get started")
			if err != nil {
				log.Printf("could not send guild creation message: %s", err)
			}
			log.Printf("Guild %s greeted", event.Name)
			return
		}
	}
}
