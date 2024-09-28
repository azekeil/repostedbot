package reposted

import (
	"strconv"
	"time"

	"github.com/azekeil/repostedbot/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func ScoreSummary(s *discordgo.Session, m *discordgo.MessageCreate) {
	var authors, scores string
	for authorID, score := range Scores[m.GuildID] {
		authors += GetUserLink(authorID) + "\n"
		scores += strconv.Itoa(len(score)) + "\n"
	}

	msg := bot.NewEmbed("")
	msg.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Author",
			Value:  authors,
			Inline: true,
		},
		{
			Name:   "Score",
			Value:  scores,
			Inline: true,
		},
	}
	bot.SendEmbed(s, m.ChannelID, msg)
}

func ScoreDetails(s *discordgo.Session, m *discordgo.MessageCreate, authorHandle string) {
	authorID := GetAuthorIDfromLink(authorHandle)
	var timestamps, reposts, originals string
	for _, repost := range Scores[m.GuildID][authorID] {
		timestamps += repost.TimeStamp.Format(time.DateTime) + "\n"
		reposts += GetMessageLink(repost.Ref) + "\n"
		originals += GetMessageLink(repost.OriginalRef) + "\n"
	}
	msg := bot.NewEmbed(authorHandle + ": " + strconv.Itoa(len(Scores[m.GuildID][authorID])))
	msg.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Timestamp",
			Value:  timestamps,
			Inline: true,
		},
		{
			Name:   "Repost",
			Value:  reposts,
			Inline: true,
		},
		{
			Name:   "Original",
			Value:  originals,
			Inline: true,
		},
	}
	bot.SendEmbed(s, m.ChannelID, msg)
}
