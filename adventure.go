package goadventure

import (
	"fmt"
	"github.com/kurrik/twittergo"
	"sync"
)

func Run(useLiveTwitterClient bool) {
	var (
		twitterWrapper TwitterWrapper
		game           *Game
	)

	// set up game world
	game = &Game{}

	// set up twitter client for adventure user
	if useLiveTwitterClient {
		twitterWrapper = NewRealTwitterWrapper()
	} else {
		twitterWrapper = new(FakeTwitterWrapper)
	}

	// print some debug on the user
	twitterWrapper.PrintUserDebugInfo()

	// setup channel for listen loop to tell game loop
	// about incoming tweets
	tweetChannel := make(chan *twittergo.Tweet)
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	// setup listen loop for @mentions
	go func() {
		timeline := twitterWrapper.GetUserMentionsTimeline()
		// each tweet mentioned stuff onto channel
		for _, tweet := range *timeline {
			tweetChannel <- &tweet
		}
		close(tweetChannel)
		waitGroup.Done()
	}()

	// setup gameplay loop
	go func() {
		// fetch tweet off channel
		for tweet := range tweetChannel {
			fmt.Printf("Tweet:   %v\n", tweet.Text())
			// set gamestate
			user := tweet.User()

			gameState := game.GetStateForUser(user.ScreenName())
			response := gameState.UpdateState(tweet.Text())

			// tweet at them their "room"
			twitterWrapper.SendResponseToUser(&user, response)
		}
	}()

	waitGroup.Wait()

}
