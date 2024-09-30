package reposted

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
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
	// Order keys
	keys := make([]string, 0, len(Scores[m.GuildID][authorID]))
	sc := make(map[string]*Score, len(Scores[m.GuildID][authorID]))
	for _, repost := range Scores[m.GuildID][authorID] {
		v := strconv.FormatInt(repost.TimeStamp.Unix(), 10) + repost.Ref.MessageID
		keys = append(keys, v)
		sc[v] = repost
	}
	sort.Strings(keys)

	var timestamps, reposts, originals string
	for _, k := range keys {
		repost := sc[k]
		timestamps += repost.TimeStamp.Format(time.DateTime) + "\n"
		reposts += GetMessageLink(repost.Ref) + "\n"
		originals += GetMessageLink(repost.OriginalRef) + "\n"
	}
	if len(timestamps)+len(reposts)+len(originals) < 1024 {
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
		return
	}
	// Message is too big for an embed so just send it plain
	textTable := &strings.Builder{}
	w := tabwriter.NewWriter(textTable, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Timestamp\tRepost\tOriginal")
	for _, k := range keys {
		repost := sc[k]
		fmt.Fprintln(w, repost.TimeStamp.Format(time.DateTime)+"\t"+GetMessageLink(repost.Ref)+"\t"+GetMessageLink(repost.OriginalRef))
	}
	w.Flush()
	msg := textTable.String()
	if len(msg) < 2000 {
		bot.SendMessage(s, m.ChannelID, msg)
		return
	}
	// Message is over 2000 characters, split it into multiple messages :/
	var c string
	for _, l := range strings.SplitAfter(msg, "\n") {
		if len(c)+len(l) > 2000 {
			bot.SendMessage(s, m.ChannelID, c)
			c = l
		}
		c += l
	}
	bot.SendMessage(s, m.ChannelID, c)
}
