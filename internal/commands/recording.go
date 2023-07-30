package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/azekeil/grec/internal/bot"
	"github.com/azekeil/grec/internal/self"
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
			if _, ok := bot.ActiveRecordings[m.Author.ID]; ok {
				bot.SendEmbed(s, ch.ID, bot.NewErrorEmbed(
					"Error: You already have an active recording!\nStop it first with `!grec recordstop`",
				))
				return
			}

			// OK, join the voice channel and start recording
			v, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, true, false)
			if err != nil {
				msg := fmt.Sprint("Error joining voice channel: ", err)
				log.Println(msg)
				bot.SendEmbed(s, ch.ID, bot.NewErrorEmbed(msg))
				return
			}

			// Ensure the channels are nil'd so they get recreated if being reused on closed channels.
			v.Lock()
			v.OpusRecv = nil
			v.OpusSend = nil
			v.Unlock()

			bot.StartRecording(s, m, g, vs, v)
			bot.SendEmbed(s, ch.ID, bot.NewEmbed("OK, Recording started!"))
			return
		}
	}
	log.Printf("error: could not find user %s in voice channel", m.Author.ID)
	bot.SendEmbed(s, ch.ID, bot.NewErrorEmbed(
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
	if r, ok := bot.ActiveRecordings[m.Author.ID]; ok {
		msg := fmt.Sprintf("finished recording: %s elapsed", time.Since(r.StartTime))
		bot.StopRecording(s, m.Author.ID)
		log.Println(msg)
		bot.SendEmbed(s, m.ChannelID, bot.NewEmbed(msg))
	} else {
		log.Printf("error: user %s has no active recording", m.Author.ID)
		bot.SendEmbed(s, m.ChannelID, bot.NewErrorEmbed(
			"Error: You have no active recordings!\nStart one with `!grec recordhere`",
		))
		return
	}
}
