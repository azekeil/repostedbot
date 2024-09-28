package reposted

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/corona10/goimagehash"
)

type Post struct {
	MessageReference *discordgo.MessageReference
	TimeStamp        *time.Time
	AuthorID         string
}

type ImgHashPost = map[uint64]*Post
type ScorePosts = map[string][]*Post
type LastPosts = map[string]string

var ImgHashes = map[string]ImgHashPost{}
var Scores = map[string]ScorePosts{}
var LastPost = map[string]LastPosts{}

func hashImageFromURL(url string) (*goimagehash.ImageHash, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("got %d downloading %s", res.StatusCode, url)
	}
	defer res.Body.Close()
	m, _, err := image.Decode(res.Body)
	if err != nil {
		return nil, err
	}
	return goimagehash.AverageHash(m)
}

func findRepost(hashMap ImgHashPost, hash *goimagehash.ImageHash, distance int) (*goimagehash.ImageHash, error) {
	for h := range hashMap {
		loopHash := goimagehash.NewImageHash(h, goimagehash.AHash)
		d, err := hash.Distance(loopHash)
		if err != nil {
			return nil, err
		}
		if d <= distance {
			return loopHash, nil
		}
	}
	return nil, nil
}

func HandleMessageAttachments(s *discordgo.Session, m *discordgo.MessageCreate) (msg string, msgErr bool) {
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
			log.Printf("Added %d to hashes. Now have %d hashes for guild %s.", imgHash.GetHash(), len(thisImgHashes), m.GuildID)
		}
	}
	// Update LastPost
	if LastPost[m.GuildID] == nil {
		LastPost[m.GuildID] = LastPosts{}
	}
	LastPost[m.GuildID][m.ChannelID] = m.ID
	err := SaveDB()
	if err != nil {
		log.Fatalf("Fatal error saving DB: %v", err)
	}
	return
}
