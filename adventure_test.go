package goadventure

import (
	"github.com/kurrik/twittergo"
	"testing"
)

type TestHarnessTwitterWrapper struct {
	timeToFinish chan bool
	sentMessages []string
}

func (twitterWrapper *TestHarnessTwitterWrapper) PrintUserDebugInfo() {
	// no op
}

func (twitterWrapper *TestHarnessTwitterWrapper) GetUserMentionsTimeline() *twittergo.Timeline {
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

func (twitterWrapper *TestHarnessTwitterWrapper) RespondToTweet(tweet *twittergo.Tweet, message string) {
	twitterWrapper.sentMessages = append(twitterWrapper.sentMessages, message)
	twitterWrapper.timeToFinish <- true
}

func TestRun(t *testing.T) {
	twitterWrapper := new(TestHarnessTwitterWrapper)
	twitterWrapper.timeToFinish = make(chan bool)

	Run(twitterWrapper.timeToFinish, twitterWrapper)

	if len(twitterWrapper.sentMessages) != 1 {
		t.Fatalf("Expected 1 sent twitter message got %v", len(twitterWrapper.sentMessages))
	}
}
