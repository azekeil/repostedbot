package commands

import (
	"github.com/alex-broad/grec/internal/actions"
	"github.com/alex-broad/grec/internal/self"
	"github.com/bwmarrin/discordgo"
)

// help: grec bot for recording things
//
// To see a list and summary of commands, type `!grec list`
// To see help for a specific command, type `!grec help <command>`
func (c *Command) Help(s *discordgo.Session, m *discordgo.MessageCreate, help self.DocFuncs) {
	// Send this function comment as help text
	actions.SendEmbed(s, m.ChannelID, &discordgo.MessageEmbed{
		Color:       0x1c1c1c,
		Description: help.CommandHelp("Help"),
	})
}
