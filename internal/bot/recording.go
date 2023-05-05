package bot

import (
	"log"
	"time"

	"github.com/azekeil/grec/internal/utils"
	"github.com/azekeil/grec/internal/voice"
	"github.com/bwmarrin/discordgo"
)

type ActiveRecording struct {
	User            *discordgo.User
	Guild           *discordgo.Guild
	VoiceState      *discordgo.VoiceState
	VoiceConnection *discordgo.VoiceConnection
	StartTime       time.Time
}

var ActiveRecordings = make(map[string]*ActiveRecording)

// StartRecording does not error check the pointers, so you better have.
func StartRecording(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
	g *discordgo.Guild,
	vs *discordgo.VoiceState,
	v *discordgo.VoiceConnection,
) {
	ActiveRecordings[m.Author.ID] = &ActiveRecording{
		User:            m.Author,
		Guild:           g,
		VoiceState:      vs,
		VoiceConnection: v,
		StartTime:       time.Now(),
	}
	go voice.HandleVoice(v.OpusRecv, utils.MustGetwd())
	log.Println("recording started for user", m.Author.ID)
}

func StopRecording(s *discordgo.Session, userID string) {
	if r, ok := ActiveRecordings[userID]; ok {
		func() {
			defer func() {
				res := recover()
				if _, ok := res.(error); ok {
					log.Println("Recovered from closing already closed go channel! (Suspect no active speakers)")
				}
			}()
			close(r.VoiceConnection.OpusRecv)
		}()
		r.VoiceConnection.Close()
		// This was not obvious to get the bot to quit the voice channel :(
		err := s.ChannelVoiceJoinManual(r.Guild.ID, "", true, false)
		if err != nil {
			log.Println("Error leaving voice channel:", err)
		}
		delete(ActiveRecordings, userID)
	}
}

func StopAllRecordings(s *discordgo.Session) {
	for u := range ActiveRecordings {
		StopRecording(s, u)
	}
}
