package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/alex-broad/grec/internal/actions"
	"github.com/alex-broad/grec/internal/self"
	"github.com/bwmarrin/discordgo"
)

// recordhere: join your voice channel and start recording
//
// This will attempt to save an autonamed file in the configured location
// Any errors will be fed back to the user
//
// To stop recording, use `!grec recordstop`
func (c *Command) RecordHere(s *discordgo.Session, m *discordgo.MessageCreate, help self.DocFuncs) {
	// Find the channel that the message came from.
	ch, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("error: could not find channel:", err)
		return
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(ch.GuildID)
	if err != nil {
		log.Println("error: could not find guild:", err)
		return
	}

	// Look for the message sender in that guild's current voice states.
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {

			// Check they don't already have an active recording
			if _, ok := actions.ActiveRecordings[m.Author.ID]; ok {
				actions.SendEmbed(s, ch.ID, actions.NewErrorEmbed(
					"Error: You already have an active recording!\nStop it first with `!grec recordstop`",
				))
				return
			}

			// OK, join the voice channel and start recording
			v, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, true, false)
			if err != nil {
				msg := fmt.Sprint("Error joining voice channel: ", err)
				log.Println(msg)
				actions.SendEmbed(s, ch.ID, actions.NewErrorEmbed(msg))
				return
			}

			actions.StartRecording(s, m, g, vs, v)
			actions.SendEmbed(s, ch.ID, actions.NewEmbed("OK, Recording started!"))
			return
		}
	}
	log.Printf("error: could not find user %s in voice channel", m.Author.ID)
	actions.SendEmbed(s, ch.ID, actions.NewErrorEmbed(
		"Error: Can't find you in a voice channel.\nPlease ensure you are in a voice channel before starting recording!",
	))
}

// recordstop: stop any ongoing recording
//
// This will stop any ongoing recording and return the filename and some stats
// Any errors will be fed back to the user
//
// To start recording, use `!grec recordstart`
func (c *Command) RecordStop(s *discordgo.Session, m *discordgo.MessageCreate, help self.DocFuncs) {

	// See if we have an active recording for this user
	if r, ok := actions.ActiveRecordings[m.Author.ID]; ok {
		msg := fmt.Sprintf("finished recording: %s elapsed", time.Since(r.StartTime))
		actions.StopRecording(s, m.Author.ID)
		log.Println(msg)
		actions.SendEmbed(s, m.ChannelID, actions.NewEmbed(msg))
	} else {
		log.Printf("error: user %s has no active recording", m.Author.ID)
		actions.SendEmbed(s, m.ChannelID, actions.NewErrorEmbed(
			"Error: You have no active recordings!\nStart one with `!grec recordhere`",
		))
		return
	}
}
