package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func SendMessage(s *discordgo.Session, channelID, content string) {
	_, err := s.ChannelMessageSend(channelID, content)
	if err != nil {
		log.Printf("could not send message: %s", err)
	}
}

func NewEmbed(description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Description: description,
		Color:       0x1c1c1c,
	}
}

func NewErrorEmbed(description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Description: description,
		Color:       0xb40000,
	}
}

func SendEmbed(s *discordgo.Session, channelID string, embed *discordgo.MessageEmbed) {
	_, err := s.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		log.Printf("could not send embed: %s", err)
	}
}
