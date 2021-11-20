package actions

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

func SendEmbed(s *discordgo.Session, channelID string, embed *discordgo.MessageEmbed) {
	_, err := s.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		log.Printf("could not send embed: %s", err)
	}
}
