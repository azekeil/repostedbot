package reposted

import (
	"fmt"
	"log"
	"net/url"

	"github.com/bwmarrin/discordgo"
	"github.com/corona10/goimagehash"
)

type guild struct {
	i ImgHashPost
	s ScorePosts
	l LastPost
}

func newGuild(guildID string) *guild {
	// Ensure this guild's hashes are initialized
	if ImgHashes.Get(guildID) == nil {
		ImgHashes.Set(guildID, NewSafeMap[uint64, *Post]())
	}
	if Scores.Get(guildID) == nil {
		Scores.Set(guildID, NewSafeMap[string, []*Score]())
	}
	if LastPosts.Get(guildID) == nil {
		LastPosts.Set(guildID, NewSafeMap[string, string]())
	}
	return &guild{
		i: ImgHashes.Get(guildID),
		s: Scores.Get(guildID),
		l: LastPosts.Get(guildID),
	}
}

// MessageHandler is the handler for incoming messages.
func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) (msg string, msgErr bool) {
	g := newGuild(m.GuildID)
	msg = g.processMessage(s, m.Message, msg)
	if err := SaveDB(); err != nil {
		log.Fatalf("Fatal error saving DB: %v", err)
	}
	return msg, false
}

func (g *guild) processMessage(s *discordgo.Session, m *discordgo.Message, msg string) string {
	// If the content is just a URL, process it.
	u, err := url.ParseRequestURI(m.Content)
	if err == nil {
		msg = g.processImage(s, m, 0, u.String(), msg)
	}
	// Now do any attachments.
	for i, a := range m.Attachments {
		msg = g.processImage(s, m, i, a.URL, msg)
	}
	// Update LastPost
	g.updateLastPost(m)
	return msg
}

func (g *guild) processImage(s *discordgo.Session, m *discordgo.Message, attachmentNum int, URL, msg string) string {
	imgHash, original := g.hashURL(m, URL)
	if original != nil {
		// Repost found! Add to score
		g.addScore(m, original)
		// And add something to the message to return.
		msg += g.addRepostToMsg(m, attachmentNum, msg, original)
	}
	if imgHash != nil {
		// Now add post to DB
		g.addToDB(s, m, imgHash)
	}
	return msg
}

func (g *guild) updateLastPost(m *discordgo.Message) {
	if m.ID > g.l.Get(m.ChannelID) {
		g.l.Set(m.ChannelID, m.ID)
	}
}

func (g *guild) addScore(m *discordgo.Message, original *goimagehash.ImageHash) {
	originalImgHash := g.i.Get(original.GetHash())
	g.s.Set(m.Author.ID,  append(g.s.Get(m.Author.ID), &Score{
		Ref:              m.Reference(),
		TimeStamp:        &m.Timestamp,
		AuthorID:         m.Author.ID,
		OriginalRef:      originalImgHash.MessageReference,
		OriginalAuthorID: originalImgHash.AuthorID,
	}))
}

func (g *guild) addToDB(s *discordgo.Session, m *discordgo.Message, imgHash *goimagehash.ImageHash) {
	if _, ok := g.i.Get2(imgHash.GetHash()); !ok {
		g.i.Set(imgHash.GetHash(), &Post{
			MessageReference: m.Reference(),
			AuthorID:         m.Author.ID,
		})
		log.Printf("Added %d to hashes. Now have %d hashes for guild %s.", imgHash.GetHash(), g.i.Len(), GetGuildName(s, m.GuildID))
	}
}

func (g *guild) addRepostToMsg(m *discordgo.Message, i int, msg string, repost *goimagehash.ImageHash) string {
	aNumStr := ""
	if len(m.Attachments) > 1 {
		aNumStr = fmt.Sprintf("image %d/%d is a ", i, len(m.Attachments))
	}
	msg += fmt.Sprintf("%srepost of %s by %s! That's %d reposts %s has made now ;)\n",
		aNumStr,
		GetMessageLink(g.i.Get(repost.GetHash()).MessageReference),
		GetUserLink(g.i.Get(repost.GetHash()).AuthorID),
		len(g.s.Get(m.Author.ID)),
		GetUserLink(m.Author.ID),
	)
	return msg
}

// hashURL downloads and generates a hash for an attachment.
// If it's a repost then repost will contain the hash of the repost, else nil.
func (g *guild) hashURL(m *discordgo.Message, URL string) (imgHash *goimagehash.ImageHash, original *goimagehash.ImageHash) {
	imgHash, err := hashImageFromURL(URL)
	if err != nil {
		log.Printf("failed to process %s: %v", m.ID, err)
		return
	}
	original, err = findRepost(g.i, imgHash, 2)
	if err != nil {
		log.Printf("failed to findRepost: %v", err)
	}
	return
}
