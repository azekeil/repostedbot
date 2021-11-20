package commands

import (
	"fmt"

	"github.com/alex-broad/grec/internal/actions"
	"github.com/alex-broad/grec/internal/self"
	"github.com/bwmarrin/discordgo"
	"github.com/ryanuber/columnize"
)

type Command struct{}

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

// list: lists available commands with summaries
func (c *Command) List(s *discordgo.Session, m *discordgo.MessageCreate, help self.DocFuncs) {
	// Send all summaries as an embed
	pretty := columnize.Format(help.AllSummaries(), &columnize.Config{Delim: ":"})
	actions.SendEmbed(s, m.ChannelID, &discordgo.MessageEmbed{
		Color:       0x1c1c1c,
		Description: fmt.Sprintf("```%s```", pretty),
	})
}
