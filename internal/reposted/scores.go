package reposted

import (
	"strconv"

	"github.com/azekeil/repostedbot/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func ScoreSummary(s *discordgo.Session, m *discordgo.MessageCreate) {
	var authors, scores string
	for authorID, score := range Scores {
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
	var timestamps, posts string
	for _, post := range Scores[authorID] {
		timestamps += post.TimeStamp.String() + "\n"
		posts += GetMessageLink(post.MessageReference) + "\n"
	}
	msg := bot.NewEmbed(authorHandle + ": " + strconv.Itoa(len(Scores[authorID])))
	msg.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Timestamp",
			Value:  timestamps,
			Inline: true,
		},
		{
			Name:   "Repost",
			Value:  posts,
			Inline: true,
		},
	}
	bot.SendEmbed(s, m.ChannelID, msg)
}
