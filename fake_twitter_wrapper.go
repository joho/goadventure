package goadventure

import (
	"bufio"
	"fmt"
	"github.com/kurrik/twittergo"
	"os"
	"strconv"
	"strings"
	"time"
)

type FakeTwitterWrapper struct {
	currentTweetId uint64
}

func NewFakeTwitterWrapper() TwitterWrapper {
	return &FakeTwitterWrapper{
		uint64(time.Now().Unix()),
	}
}

func (tw *FakeTwitterWrapper) DurationUntilNextRead() time.Duration {
	return 100 * time.Millisecond
}

func (tw *FakeTwitterWrapper) GetUserMentionsTimeline(tweetChannel chan *twittergo.Tweet) {
	fmt.Print("Text to tweet @goadventure (q to quit): ")
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		os.Exit(1)
	}
	input = strings.TrimSpace(input)
	if input == "q" {
		close(tweetChannel)
	} else {
		tweet := tw.createFakeTweetForText(input)

		fmt.Printf("Hypothetically Receive tweet '%v' from '%v'\n", tweet.Text(), tweet.User().ScreenName())

		tweetChannel <- &tweet
	}
}

func (tw *FakeTwitterWrapper) RespondToTweet(tweet *twittergo.Tweet, message string) {
	fmt.Printf("Hypothetically Send tweet '%v' to '%v'\n", message, tweet.User().ScreenName())
}

func (tw *FakeTwitterWrapper) createFakeTweetForText(tweetText string) twittergo.Tweet {
	tw.currentTweetId += 1
	user := map[string]interface{}{
		"screen_name": "johnbarton",
		"id_str":      "123549854887",
	}
	tweet := twittergo.Tweet{
		"text":   "@gotextadventure " + tweetText,
		"user":   user,
		"id_str": strconv.FormatUint(tw.currentTweetId, 10),
	}
	return tweet
}
