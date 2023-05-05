package commands

import (
	"strings"

	"github.com/azekeil/grec/internal/bot"
	"github.com/azekeil/grec/internal/self"
	"github.com/bwmarrin/discordgo"
)

type Command struct{}

// help: grec bot for recording things
//
// To see a list and summary of commands, type `!grec list`
// To see help for a specific command, type `!grec help <command>`
func (c *Command) Help(s *discordgo.Session, m *discordgo.MessageCreate, help self.DocFuncs) {
	// Send this function comment as help text
	bot.SendEmbed(s, m.ChannelID, bot.NewEmbed(help.CommandHelp("Help")))
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
	msg := bot.NewEmbed("")
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
	bot.SendEmbed(s, m.ChannelID, msg)
}
