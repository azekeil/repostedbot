package reposted

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// CatchUp with the last 100 (max) messages since the last processed post in each channel
func CatchUp(s *discordgo.Session) error {
	for guildID, lpGuild := range LastPosts.Iter() {
		g := newGuild(guildID)
		for channelID, messageID := range lpGuild.Iter() {
			messages, err := s.ChannelMessages(channelID, 100, "", messageID, "")
			if err != nil {
				return err
			}
			g.bulkProcess(s, messages)
		}
	}
	return nil
}

// bulkProcess processes a list of messages into the database and updates the lastpost
func (g *guild) bulkProcess(s *discordgo.Session, messages []*discordgo.Message) {
	for _, m := range messages {
		// Workaround: Ensure message has a GuildID
		EnsureGuildID(s, m)
		g.processMessage(s, m, "")
	}
	if err := SaveDB(); err != nil {
		log.Fatalf("Fatal error saving DB: %v", err)
	}
}

// AddHistory is the underlying implementation of the addhistory command
func AddHistory(s *discordgo.Session, guildID string, messages []*discordgo.Message) {
	g := newGuild(guildID)
	g.bulkProcess(s, messages)
}
