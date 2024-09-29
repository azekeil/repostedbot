package reposted

import (
	"fmt"
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

func GetMessageLink(m *discordgo.MessageReference) string {
	return fmt.Sprintf("https://discord.com/channels/%s/%s/%s", m.GuildID, m.ChannelID, m.MessageID)
}

var refFromLinkRegex = regexp.MustCompile(
	`https:\/\/discord\.com\/channels\/(\d+)\/(\d+)\/(\d+)`)

func GetRefFromMessageLink(msgLink string) *discordgo.MessageReference {
	a := refFromLinkRegex.FindStringSubmatch(msgLink)
	if a == nil {
		return nil
	}
	b := true
	return &discordgo.MessageReference{
		MessageID:       a[2],
		ChannelID:       a[1],
		GuildID:         a[0],
		FailIfNotExists: &b,
	}
}

func GetUserLink(ID string) string {
	return fmt.Sprintf("<@%s>", ID)
}

func GetAuthorIDfromLink(link string) string {
	return link[2 : len(link)-1]
}

func GetGuildName(s *discordgo.Session, guildID string) string {
	guild, err := s.Guild(guildID)
	if err != nil {
		log.Printf("failed to get Guild: %v", err)
		return guildID
	}
	return guild.Name
}
