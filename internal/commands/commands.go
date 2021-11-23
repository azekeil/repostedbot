package commands

import (
	"strings"

	"github.com/alex-broad/grec/internal/actions"
	"github.com/alex-broad/grec/internal/self"
	"github.com/bwmarrin/discordgo"
)

type Command struct{}

// help: grec bot for recording things
//
// To see a list and summary of commands, type `!grec list`
// To see help for a specific command, type `!grec help <command>`
func (c *Command) Help(s *discordgo.Session, m *discordgo.MessageCreate, help self.DocFuncs) {
	// Send this function comment as help text
	actions.SendEmbed(s, m.ChannelID, actions.NewEmbed(help.CommandHelp("Help")))
}

// list: lists available commands with summaries
func (c *Command) List(s *discordgo.Session, m *discordgo.MessageCreate, help self.DocFuncs) {
	// Send all summaries as an embed
	var cmd, sum string
	for _, l := range help.AllSummaries() {
		sp := strings.SplitN(l, ":", 2)
		cmd += sp[0] + "\n"
		sum += sp[1] + "\n"
	}
	msg := actions.NewEmbed("")
	msg.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Command",
			Value:  cmd,
			Inline: true,
		},
		{
			Name:   "Summary",
			Value:  sum,
			Inline: true,
		},
	}
	actions.SendEmbed(s, m.ChannelID, msg)
}
