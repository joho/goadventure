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

func (tw *FakeTwitterWrapper) GetUserMentionsTimeline() *twittergo.Timeline {
	tweet := createFakeTweetForText("Go North")
	fmt.Printf("Hypothetically Receive tweet '%v' from '%v'\n", tweet.Text(), tweet.User().ScreenName())
	return &twittergo.Timeline{tweet}
}

func (tw *FakeTwitterWrapper) RespondToTweet(tweet *twittergo.Tweet, message string) {
	fmt.Printf("Hypothetically Send tweet '%v' to '%v'\n", message, tweet.User().ScreenName())
}

func createFakeTweetForText(tweetText string) twittergo.Tweet {
	user := map[string]interface{}{
		"screen_name": "johnbarton",
		"id_str":      "123549854887",
	}
	tweet := twittergo.Tweet{
		"text": "@gotextadventure " + tweetText,
		"user": user,
	}
	return tweet
}
