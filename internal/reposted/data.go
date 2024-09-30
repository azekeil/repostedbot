package reposted

import (
	"time"

	"github.com/bwmarrin/discordgo"
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
	ImgHashPost = *SafeMap[uint64, *Post]
	ScorePosts  = *SafeMap[string, []*Score]
	LastPost    = *SafeMap[string, string]
)

var (
	ImgHashes = SafeMap[string, ImgHashPost]{}
	Scores    = SafeMap[string, ScorePosts]{}
	LastPosts = SafeMap[string, LastPost]{}
)
