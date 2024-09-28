package reposted

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/corona10/goimagehash"
)

type Post struct {
	MessageReference *discordgo.MessageReference
	AuthorID         string
}

type Score struct {
	Ref              *discordgo.MessageReference
	TimeStamp        *time.Time
	AuthorID         string
	OriginalRef      *discordgo.MessageReference
	OriginalAuthorID string
}

type (
	ImgHashPost = map[uint64]*Post
	ScorePosts  = map[string][]*Score
	LastPost    = map[string]string
)

var (
	ImgHashes = map[string]ImgHashPost{}
	Scores    = map[string]ScorePosts{}
	LastPosts = map[string]LastPost{}
)

type guild struct {
	i ImgHashPost
	s ScorePosts
	l LastPost
}

func newGuild(guildID string) *guild {
	// Ensure this guild's hashes are initialized
	if ImgHashes[guildID] == nil {
		ImgHashes[guildID] = ImgHashPost{}
	}
	if Scores[guildID] == nil {
		Scores[guildID] = ScorePosts{}
	}
	if LastPosts[guildID] == nil {
		LastPosts[guildID] = LastPost{}
	}
	return &guild{
		i: ImgHashes[guildID],
		s: Scores[guildID],
		l: LastPosts[guildID],
	}
}

func HandleMessageAttachments(s *discordgo.Session, m *discordgo.MessageCreate) (msg string, msgErr bool) {
	g := newGuild(m.GuildID)
	for i, a := range m.Attachments {
		imgHash, original := g.processAttachment(m.Message, a.URL)
		if original != nil {
			// Repost found! Add to score
			g.addScore(m.Message, original)
			// And add something to the message to return.
			msg = g.addRepostToMsg(m.Message, i, msg, original)
		}
		// Now add post to DB
		g.addToDB(s, m.Message, imgHash)
	}
	// Update LastPost
	g.updateLastPost(m.Message)
	err := SaveDB()
	if err != nil {
		log.Fatalf("Fatal error saving DB: %v", err)
	}
	return
}

func (g *guild) updateLastPost(m *discordgo.Message) {
	g.l[m.ChannelID] = m.ID
}

func (g *guild) addScore(m *discordgo.Message, original *goimagehash.ImageHash) {
	originalImgHash := g.i[original.GetHash()]
	g.s[m.Author.ID] = append(g.s[m.Author.ID], &Score{
		Ref:              m.Reference(),
		TimeStamp:        &m.Timestamp,
		AuthorID:         m.Author.ID,
		OriginalRef:      originalImgHash.MessageReference,
		OriginalAuthorID: originalImgHash.AuthorID,
	})
}

func (g *guild) addToDB(s *discordgo.Session, m *discordgo.Message, imgHash *goimagehash.ImageHash) {
	if _, ok := g.i[imgHash.GetHash()]; !ok {
		g.i[imgHash.GetHash()] = &Post{
			MessageReference: m.Reference(),
			AuthorID:         m.Author.ID,
		}
		log.Printf("Added %d to hashes. Now have %d hashes for guild %s.", imgHash.GetHash(), len(g.i), GetGuildName(s, m.GuildID))
	}
}

func (g *guild) addRepostToMsg(m *discordgo.Message, i int, msg string, repost *goimagehash.ImageHash) string {
	aNumStr := ""
	if len(m.Attachments) > 1 {
		aNumStr = fmt.Sprintf("image %d/%d is a ", i, len(m.Attachments))
	}
	msg += fmt.Sprintf("%srepost of %s by %s! That's %d reposts %s has made now ;)\n",
		aNumStr,
		GetMessageLink(g.i[repost.GetHash()].MessageReference),
		GetUserLink(g.i[repost.GetHash()].AuthorID),
		len(g.s[m.Author.ID]),
		GetUserLink(m.Author.ID),
	)
	return msg
}

// processAttachment downloads and generates a hash for an attachment.
// If it's a repost then repost will contain the hash of the repost, else nil.
func (g *guild) processAttachment(m *discordgo.Message, url string) (imgHash *goimagehash.ImageHash, original *goimagehash.ImageHash) {
	imgHash, err := hashImageFromURL(url)
	if err != nil {
		log.Printf("failed to process %s: %v", m.ID, err)
	}
	original, err = findRepost(g.i, imgHash, 2)
	if err != nil {
		log.Printf("failed to findRepost: %v", err)
	}
	return imgHash, original
}
