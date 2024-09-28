package reposted

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Post struct {
	MessageReference *discordgo.MessageReference
	TimeStamp        *time.Time
	AuthorID         string
}

type ImgHashPost = map[uint64]*Post
type ScorePosts = map[string][]*Post
type LastPost = map[string]string

var ImgHashes = map[string]ImgHashPost{}
var Scores = map[string]ScorePosts{}
var LastPosts = map[string]LastPost{}

// func ProcessMessage(m *discordgo.Message) (reposts []*goimagehash.ImageHash, err error) {
// 	for _, a := range m.Attachments {
// 		repost, err := ProcessAttachment(a)
// 		if err != nil {
// 			return nil, err
// 		}
// 		reposts = append(reposts, repost)
// 	}
// 	return reposts, nil
// }

// func ProcessAttachment(a *discordgo.MessageAttachment) (repost *goimagehash.ImageHash, err error) {
// 	imgHash, err := hashImageFromURL(a.URL)
// 	if err != nil {
// 		log.Printf("failed to process %s: %v", m.Message.ID, err)
// 	}
// 	repost, err := findRepost(thisImgHashes, imgHash, 2)
// 	if err != nil {
// 		log.Printf("failed to findRepost: %v", err)
// 	}
// }

func HandleMessageAttachments(s *discordgo.Session, m *discordgo.MessageCreate) (msg string, msgErr bool) {
	// Ensure this guild's hashes are initialized
	if ImgHashes[m.GuildID] == nil {
		ImgHashes[m.GuildID] = ImgHashPost{}
	}
	thisImgHashes := ImgHashes[m.GuildID]
	if Scores[m.GuildID] == nil {
		Scores[m.GuildID] = ScorePosts{}
	}
	thisScores := Scores[m.GuildID]
	for i, a := range m.Attachments {
		imgHash, err := hashImageFromURL(a.URL)
		if err != nil {
			log.Printf("failed to process %s: %v", m.Message.ID, err)
		}
		repost, err := findRepost(thisImgHashes, imgHash, 2)
		if err != nil {
			log.Printf("failed to findRepost: %v", err)
		}
		if repost != nil {
			// Repost found! Add to score
			thisScores[m.Author.ID] = append(thisScores[m.Author.ID], &Post{
				MessageReference: m.Reference(),
				TimeStamp:        &m.Timestamp,
				AuthorID:         m.Author.ID,
			})
			// And add something to the message to return.
			aNumStr := ""
			if len(m.Attachments) > 1 {
				aNumStr = fmt.Sprintf("image %d/%d is a ", i, len(m.Attachments))
			}
			msg += fmt.Sprintf("%srepost of %s by %s! That's %d reposts %s has made now ;)\n",
				aNumStr,
				GetMessageLink(thisImgHashes[repost.GetHash()].MessageReference),
				GetUserLink(thisImgHashes[repost.GetHash()].AuthorID),
				len(thisScores[m.Author.ID]),
				GetUserLink(m.Author.ID),
			)
		}
		// Now add post to DB
		if _, ok := thisImgHashes[imgHash.GetHash()]; !ok {
			thisImgHashes[imgHash.GetHash()] = &Post{
				MessageReference: m.Reference(),
				TimeStamp:        &m.Timestamp,
				AuthorID:         m.Author.ID,
			}
			var guildName string
			guild, err := s.Guild(m.GuildID)
			if err != nil {
				log.Printf("failed to get Guild: %v", err)
			} else {
				guildName = guild.Name
			}
			log.Printf("Added %d to hashes. Now have %d hashes for guild %s.", imgHash.GetHash(), len(thisImgHashes), guildName)
		}
	}
	// Update LastPost
	if LastPosts[m.GuildID] == nil {
		LastPosts[m.GuildID] = LastPost{}
	}
	LastPosts[m.GuildID][m.ChannelID] = m.ID
	err := SaveDB()
	if err != nil {
		log.Fatalf("Fatal error saving DB: %v", err)
	}
	return
}
