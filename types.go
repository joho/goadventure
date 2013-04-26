package goadventure

import (
	"github.com/kurrik/twittergo"
	"time"
)

type TwitterWrapper interface {
	DurationUntilNextRead() time.Duration
	GetUserMentionsTimeline() *twittergo.Timeline
	RespondToTweet(*twittergo.Tweet, string)
}
