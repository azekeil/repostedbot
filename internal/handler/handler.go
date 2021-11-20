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

		// check if the message starts with "!grec"
		if strings.HasPrefix(m.Content, "!grec") {

			all := strings.Split(m.Content, " ")
			if len(all) < 2 {
				c.Help(s, m, help)
				return
			}

			// OK, we know there must be a command..
			command := strings.Title(all[1])

			// See if there's a function with the same name, if so, call it
			var msg string
			if _, ok := help[command]; ok {
				args := make(map[string]interface{}, 0)
				args["s"] = s
				args["m"] = m
				args["help"] = help
				self.CallMethod(c, command, args)
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
}
