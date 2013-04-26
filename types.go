package goadventure

import "github.com/kurrik/twittergo"

type TwitterWrapper interface {
	GetUserMentionsTimeline() *twittergo.Timeline
	RespondToTweet(*twittergo.Tweet, string)
}
