package goadventure

import (
	"github.com/kurrik/twittergo"
	"testing"
)

type TestHarnessTwitterWrapper struct {
	sentMessages []string
}

func (twitterWrapper *TestHarnessTwitterWrapper) PrintUserDebugInfo() {
	// no op
}

func (twitterWrapper *TestHarnessTwitterWrapper) GetUserMentionsTimeline() *twittergo.Timeline {
	user := map[string]interface{}{
		"screen_name": "johnbarton",
	}
	tweet := twittergo.Tweet{
		"text": "@gotextadventure go north",
		"user": user,
	}
	return &twittergo.Timeline{tweet}
}

func (twitterWrapper *TestHarnessTwitterWrapper) SendResponseToUser(user *twittergo.User, message string) {
	twitterWrapper.sentMessages = append(twitterWrapper.sentMessages, message)
}

func TestRun(t *testing.T) {
	twitterWrapper := new(TestHarnessTwitterWrapper)
	Run(twitterWrapper)
	if len(twitterWrapper.sentMessages) != 1 {
		t.Fatalf("Expected 1 sent twitter message got %v", len(twitterWrapper.sentMessages))
	}
}
