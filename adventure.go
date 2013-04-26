package goadventure

import (
	"github.com/kurrik/twittergo"
	"sync"
	"time"
)

func Run(stopRunning chan bool, twitterWrapper TwitterWrapper) {
	var (
		game      *Game
		tweetRepo TweetRepo
	)

	// set up game world
	game = CreateGame()
	tweetRepo = CreateTweetLogger()

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
				close(tweetChannel)
				break ListenLoop
			default:
				if time.Since(timelineLastReadAt) > twitterWrapper.DurationUntilNextRead() {
					timeline := twitterWrapper.GetUserMentionsTimeline()
					timelineLastReadAt = time.Now()

					// each tweet mentioned stuff onto channel
					for _, tweet := range *timeline {
						tweetChannel <- &tweet
					}
				}
			}
		}

	}()

	// setup gameplay loop
	go func() {

		// fetch tweet off channel
		for tweet := range tweetChannel {
			if !tweetRepo.TweetAlreadyHandled(tweet.Id()) {
				// play the game
				message := game.Play(tweet.User().Id(), tweet.Text())

				// tweet at them their "room"
				twitterWrapper.RespondToTweet(tweet, message)
				tweetRepo.StoreTweetHandled(tweet.Id(), tweet.Text())
			}
		}

		waitGroup.Done()

	}()

	waitGroup.Wait()

}
