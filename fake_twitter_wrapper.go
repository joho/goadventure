package goadventure

import (
	"bufio"
	"fmt"
	"github.com/kurrik/twittergo"
	"os"
	"strings"
	"time"
)

type FakeTwitterWrapper struct{}

func (tw *FakeTwitterWrapper) DurationUntilNextRead() time.Duration {
	return 100 * time.Millisecond
}

func (tw *FakeTwitterWrapper) GetUserMentionsTimeline() *twittergo.Timeline {
	fmt.Print("Text to tweet @goadventure: ")
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		// fmt.Printf("Error reading tweet text: %s", err)
		os.Exit(1)
	}
	input = strings.TrimSpace(input)
	tweet := createFakeTweetForText(input)

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
