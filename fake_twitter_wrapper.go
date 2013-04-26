package goadventure

import (
	"fmt"
	"github.com/kurrik/twittergo"
	"time"
)

type FakeTwitterWrapper struct{}

func (tw *FakeTwitterWrapper) DurationUntilNextRead() time.Duration {
	return 1 * time.Second
}

func (twitterWrapper *FakeTwitterWrapper) GetUserMentionsTimeline() *twittergo.Timeline {
	user := map[string]interface{}{
		"screen_name": "johnbarton",
		"id_str":      "123549854887",
	}
	tweet := twittergo.Tweet{
		"text": "@gotextadventure go north",
		"user": user,
	}
	return &twittergo.Timeline{tweet}
}

func (twitterWrapper *FakeTwitterWrapper) RespondToTweet(tweet *twittergo.Tweet, message string) {
	fmt.Printf("Hypothetically Send tweet '%v' to '%v'\n", message, tweet.User().ScreenName())
}
