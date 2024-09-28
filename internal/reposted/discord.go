package reposted

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func GetMessageLink(m *discordgo.MessageReference) string {
	return fmt.Sprintf("https://discord.com/channels/%s/%s/%s", m.GuildID, m.ChannelID, m.MessageID)
}

func GetUserLink(ID string) string {
	return fmt.Sprintf("<@%s>", ID)
}
