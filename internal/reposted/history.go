package reposted

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// CatchUp with the last 100 (max) messages since the last processed post in each channel
func CatchUp(s *discordgo.Session) error {
	for guildID, lpGuild := range LastPosts {
		g := newGuild(guildID)
		for channelID, messageID := range lpGuild {
			messages, err := s.ChannelMessages(channelID, 100, "", messageID, "")
			if err != nil {
				return err
			}
			g.BulkProcess(s, messages)
		}
	}
	return nil
}

// BulkProcess processes a list of messages into the database and updates the lastpost
func (g *guild) BulkProcess(s *discordgo.Session, messages []*discordgo.Message) {
	for _, m := range messages {
		g.processMessage(s, m, "")
	}
	if err := SaveDB(); err != nil {
		log.Fatalf("Fatal error saving DB: %v", err)
	}
}
