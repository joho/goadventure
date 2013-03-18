package goadventure

import (
	"github.com/kurrik/twittergo"
	"sync"
)

func Run(twitterWrapper TwitterWrapper) {
	var (
		game *Game
	)

	// set up game world
	game = &Game{}

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
			// set gamestate
			user := tweet.User()

			message := game.Play(user.Id(), tweet.Text())

			// tweet at them their "room"
			twitterWrapper.RespondToTweet(tweet, message)
		}
	}()

	waitGroup.Wait()

}
