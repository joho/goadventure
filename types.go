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

type TweetRepo interface {
	TweetAlreadyHandled(uint64) bool
	StoreTweetHandled(uint64, string)
}

type GameStateRepo interface {
	SetCurrentSceneForUser(uint64, *Scene)
	GetCurrentSceneForUser(uint64) *Scene
}
