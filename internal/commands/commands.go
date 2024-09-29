package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/azekeil/repostedbot/internal/bot"
	"github.com/azekeil/repostedbot/internal/reposted"
	"github.com/azekeil/repostedbot/internal/self"
	"github.com/bwmarrin/discordgo"
)

type Command struct{}

// help: repostedbot bot for detecting reposts
//
// To see a list and summary of commands, type `!rp list`
// To see help for a specific command, type `!rp help <command>`
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

// scores: list all repost scores. To list all posts by a user append their @
func (c *Command) Scores(s *discordgo.Session, m *discordgo.MessageCreate, help self.DocFuncs) {
	all := strings.Fields(m.Content)
	if len(all) == 2 {
		reposted.ScoreSummary(s, m)
		return
	}
	reposted.ScoreDetails(s, m, all[2])
}

// addhistory: add posts to repostedbot's database for tracking.
//
// To add history for up to 100 (max) latest messages:
// !rp addhistory 100
// To add history for up to 100 (max) messages before a certain message:
// !rp addhistory 100 before <message link>
// To add history for up to 100 (max) messages after a certain message:
// !rp addhistory 100 after <message link>
func (c *Command) Addhistory(s *discordgo.Session, m *discordgo.MessageCreate, help self.DocFuncs) {
	all := strings.Fields(m.Content)
	if len(all) != 3 && len(all) != 5 {
		bot.SendEmbed(s, m.ChannelID, bot.NewErrorEmbed(
			"Command must have 3 or 5 parameters:\n\n"+
				help.CommandHelp("Addhistory")))
		return
	}
	limit, err := strconv.Atoi(all[2])
	if err != nil {
		bot.SendEmbed(s, m.ChannelID, bot.NewErrorEmbed(
			fmt.Sprintf("Error converting string to int: %v", err)))
		return
	}
	var ref *discordgo.MessageReference
	var befAfter string
	var beforeID, afterID string
	channelID := m.ChannelID
	if len(all) == 5 {
		befAfter = strings.ToLower(all[3])
		if befAfter != "before" && befAfter != "after" {
			bot.SendEmbed(s, m.ChannelID, bot.NewErrorEmbed(
				"4th parameter must be 'before' or 'after'"))
			return
		} 
		ref = reposted.GetRefFromMessageLink(all[4])
		if ref == nil {
			bot.SendEmbed(s, m.ChannelID, bot.NewErrorEmbed(
				"Unable to parse message link"))
			return
		}
		channelID = ref.ChannelID
		if befAfter == "before" {
			beforeID = ref.MessageID
		} else {
			afterID = ref.MessageID
		}
	}
	messages, err := s.ChannelMessages(channelID, limit, beforeID, afterID, "")
	if err != nil {
		bot.SendEmbed(s, m.ChannelID, bot.NewErrorEmbed(
			fmt.Sprintf("Error fetching messages: %v", err)))
		return
	}
	reposted.AddHistory(s, m.GuildID, messages)
	bot.SendEmbed(s, m.ChannelID, bot.NewEmbed("Messages added successfully"))
}
