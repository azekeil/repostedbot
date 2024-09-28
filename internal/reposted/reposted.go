package reposted

import (
	"fmt"
	"image"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/corona10/goimagehash"
)

type Post struct {
	MessageReference *discordgo.MessageReference
	AuthorID         string
}

var ImgHashes = map[*goimagehash.ImageHash]*Post{}
var Scores = map[string]uint32{}

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

func findRepost(hash *goimagehash.ImageHash, distance int) (*goimagehash.ImageHash, error) {
	for h := range ImgHashes {
		d, err := hash.Distance(h)
		if err != nil {
			return nil, err
		}
		if d <= distance {
			return h, nil
		}
	}
	return nil, nil
}

func HandleMessageAttachments(s *discordgo.Session, m *discordgo.MessageCreate) (msg string, msgErr bool) {
	for i, a := range m.Attachments {
		imgHash, err := hashImageFromURL(a.URL)
		if err != nil {
			log.Printf("failed to process %s: %v", m.Message.ID, err)
		}
		repost, err := findRepost(imgHash, 2)
		if err != nil {
			log.Printf("failed to findRepost: %v", err)
		}
		if repost != nil {
			// Repost found! Add to score
			Scores[m.Author.ID]++
			// And add something to the message to return.
			aNumStr := ""
			if len(m.Attachments) > 1 {
				aNumStr = fmt.Sprintf("image %d/%d is a ", i, len(m.Attachments))
			}
			msg += fmt.Sprintf("%srepost of %s by %s! That's %d reposts %s has made now ;)\n",
				aNumStr,
				GetMessageLink(ImgHashes[repost].MessageReference),
				GetUserLink(ImgHashes[repost].AuthorID),
				Scores[m.Author.ID],
				GetUserLink(m.Author.ID),
			)
		}
		// Now add post to DB
		if _, ok := ImgHashes[imgHash]; !ok {
			ImgHashes[imgHash] = &Post{
				MessageReference: m.Reference(),
				AuthorID:         m.Author.ID,
			}
			log.Printf("Added %d to hashes. Now have %d hashes.", imgHash.GetHash(), len(ImgHashes))
		}
	}
	return
}
