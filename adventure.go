package goadventure

import (
	"github.com/kurrik/twittergo"
	"sync"
	"time"
)

/*
	Main game loop.

	* stopRunning is the control channel, pass a "true" onto it and the game loop will halt as soon as possible

	* twitterWrapper is the input/output wrapper for the game. In dev you can pass in an interactive wrapper to drive the game from the command line

	* storageEngine is how game state & tweets are stored. Use the InMemoryStorageEngine for playing around in dev.
*/
func Run(stopRunning chan bool, twitterWrapper TwitterWrapper, storageEngine StorageEngine) {
	var (
		game *Game
	)

	// set up game world
	game = CreateGame(storageEngine)

	// setup channel for listen loop to tell game loop
	// about incoming tweets
	tweetChannel := make(chan *twittergo.Tweet)
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	// setup listen loop for @mentions
	go func() {

		var timelineLastReadAt time.Time

	ListenLoop:
		for {
			select {
			case <-stopRunning:
				break ListenLoop
			default:
				if time.Since(timelineLastReadAt) > twitterWrapper.DurationUntilNextRead() {
					twitterWrapper.GetUserMentionsTimeline(tweetChannel)
					timelineLastReadAt = time.Now()
				}
			}
		}

		close(tweetChannel)

	}()

	// setup gameplay loop
	go func() {

		// fetch tweet off channel
		for tweet := range tweetChannel {
			if !storageEngine.TweetAlreadyHandled(tweet.Id()) {
				// play the game
				message := game.Play(tweet.User().Id(), tweet.Text())

				// tweet at them their "room"
				twitterWrapper.RespondToTweet(tweet, message)
				storageEngine.StoreTweetHandled(tweet.Id(), tweet.Text())
			}
		}

		waitGroup.Done()

	}()

	waitGroup.Wait()

}
