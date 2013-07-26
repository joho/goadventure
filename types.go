package goadventure

import (
	"github.com/kurrik/twittergo"
	"time"
)

type TwitterWrapper interface {
	DurationUntilNextRead() time.Duration
	GetUserMentionsTimeline(chan *twittergo.Tweet)
	RespondToTweet(*twittergo.Tweet, string)
}

type StorageEngine interface {
	TweetAlreadyHandled(uint64) bool
	StoreTweetHandled(uint64, string)
	SetCurrentSceneKeyForUser(uint64, string)
	GetCurrentSceneKeyForUser(uint64) (string, error)
}
