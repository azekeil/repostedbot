package commands

import (
	"log"
	"time"

	"github.com/alex-broad/grec/internal/self"
	"github.com/alex-broad/grec/internal/voice"
	"github.com/bwmarrin/discordgo"
)

type ActiveRecording struct {
	User            *discordgo.User
	Guild           *discordgo.Guild
	VoiceState      *discordgo.VoiceState
	VoiceConnection *discordgo.VoiceConnection
	StartTime       time.Time
}

var ActiveRecordings = make(map[string]*ActiveRecording, 0)

// recordhere: join the voice channel you are in and start recording
//
// This will attempt to save an autonamed file in the configured location
// Any errors will be fed back to the user
//
// To stop recording, use `!grec recordstop`
func (c *Command) RecordHere(s *discordgo.Session, m *discordgo.MessageCreate, help self.DocFuncs) {
	// Find the channel that the message came from.
	ch, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("error: could not find channel")
		return
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(ch.GuildID)
	if err != nil {
		log.Println("error: could not find guild")
		return
	}

	// Look for the message sender in that guild's current voice states.
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {

			// Check they don't already have an active recording
			if _, ok := ActiveRecordings[m.Author.ID]; ok {
				log.Println("error: user already has active recording")
				return
			}

			// OK, join the voice channel and start recording
			v, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, true, false)
			if err != nil {
				log.Println("Error joining voice channel: ", err)
			}

			ActiveRecordings[m.Author.ID] = &ActiveRecording{
				User:            m.Author,
				Guild:           g,
				VoiceState:      vs,
				VoiceConnection: v,
				StartTime:       time.Now(),
			}
			go voice.HandleVoice(v.OpusRecv)
			log.Println("started recording")
			return
		}
	}
	log.Println("error: could not find user in voice channel")
}

// recordstop: stop any ongoing recording
//
// This will stop any ongoing recording and return the filename and some stats
// Any errors will be fed back to the user
//
// To start recording, use `!grec recordstart`
func (c *Command) RecordStop(s *discordgo.Session, m *discordgo.MessageCreate, help self.DocFuncs) {

	// See if we have an active recording for this user
	if r, ok := ActiveRecordings[m.Author.ID]; ok {
		close(r.VoiceConnection.OpusRecv)
		r.VoiceConnection.Close()
		// This was not obvious :(
		err := s.ChannelVoiceJoinManual(r.Guild.ID,"",true, false)
		if err != nil {
			log.Println("Error leaving voice channel: ", err)
		}
		log.Printf("finished recording: %s elapsed", time.Since(r.StartTime))
		delete(ActiveRecordings, m.Author.ID)
	} else {
		log.Println("error: user has no active recording")
		return
	}
}
