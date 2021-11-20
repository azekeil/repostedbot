package handler

import (
	"fmt"
	"strings"

	"github.com/alex-broad/grec/internal/actions"
	"github.com/alex-broad/grec/internal/commands"
	"github.com/alex-broad/grec/internal/self"
	"github.com/bwmarrin/discordgo"
)

func MakeMessageCreateHandlerFunc(help self.DocFuncs) func(*discordgo.Session, *discordgo.MessageCreate) {
	c := new(commands.Command)
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {

		// Ignore all messages created by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Ignore any messages that don't start with "!grec"
		if !strings.HasPrefix(m.Content, "!grec") {
			return
		}

		all := strings.Fields(m.Content)
		if len(all) < 2 {
			c.Help(s, m, help)
			return
		}

		// OK, we know there must be a command..
		command := strings.Title(all[1])

		// See if there's a function with the same name, if so, call it
		var msg string
		if _, ok := help[command]; ok {
			self.CallMethod(c, command, []interface{}{s, m, help})
		} else {
			msg = fmt.Sprintf("Could not find command `%s`. Try `!grec list`", all[1])
		}

		// Send message if defined
		if msg != "" {
			actions.SendEmbed(s, m.ChannelID, &discordgo.MessageEmbed{
				Color:       0x1c1c1c,
				Description: msg,
			})
		}
	}
}
