package goadventure

import "github.com/kurrik/twittergo"

type TwitterWrapper interface {
	PrintUserDebugInfo()
	GetUserMentionsTimeline() *twittergo.Timeline
	RespondToTweet(*twittergo.Tweet, string)
}
