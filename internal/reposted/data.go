package reposted

import (
	"sync"
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
	ImgHashPost = map[uint64]*Post
	ScorePosts  = map[string][]*Score
	LastPost    = map[string]string
)

var (
	ImgHashes     = map[string]ImgHashPost{}
	ImgHashesLock = sync.RWMutex{}
	Scores        = map[string]ScorePosts{}
	ScoresLock    = sync.RWMutex{}
	LastPosts     = map[string]LastPost{}
	LastPostsLock = sync.RWMutex{}
)
