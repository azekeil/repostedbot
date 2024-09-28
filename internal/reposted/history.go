package reposted

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func ProcessHistory(s *discordgo.Session) error {
	for guildID, lpGuild := range LastPosts {
		g := newGuild(guildID)
		for channelID, messageID := range lpGuild {
			messages, err := s.ChannelMessages(channelID, 100, "", messageID, "")
			if err != nil {
				return err
			}
			for _, m := range messages {
				for _, a := range m.Attachments {
					imgHash, repost := g.processAttachment(m, a.URL)
					if repost != nil {
						// Repost found! Add to score
						g.addScore(m, repost)
					}
					// Now add post to DB
					g.addToDB(s, m, imgHash)
				}
				// Update LastPost
				g.updateLastPost(m)
			}
		}
	}
	err := SaveDB()
	if err != nil {
		log.Fatalf("Fatal error saving DB: %v", err)
	}
	return nil
}
