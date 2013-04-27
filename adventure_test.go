package goadventure

import (
	"github.com/kurrik/twittergo"
	"testing"
	"time"
)

type TestHarnessTwitterWrapper struct {
	timeToFinish chan bool
	sentMessages []string
}

func (tw *TestHarnessTwitterWrapper) DurationUntilNextRead() time.Duration {
	return 1 * time.Millisecond
}

func (tw *TestHarnessTwitterWrapper) GetUserMentionsTimeline() *twittergo.Timeline {
	user := map[string]interface{}{
		"screen_name": "johnbarton",
		"id_str":      "123549854887",
	}
	tweet := twittergo.Tweet{
		"text":   "@gotextadventure go north",
		"user":   user,
		"id_str": "123543654887",
	}
	// return the same tweet twice to test duplicate handling
	timeline := &twittergo.Timeline{
		tweet,
	}
	return timeline
}

func (tw *TestHarnessTwitterWrapper) RespondToTweet(tweet *twittergo.Tweet, message string) {
	tw.sentMessages = append(tw.sentMessages, message)
	tw.timeToFinish <- true
}

func TestRun(t *testing.T) {
	tw := new(TestHarnessTwitterWrapper)
	tw.timeToFinish = make(chan bool)

	Run(tw.timeToFinish, tw)

	if len(tw.sentMessages) != 1 {
		t.Fatalf("Expected 1 sent twitter message got %v", len(tw.sentMessages))
	}
}
