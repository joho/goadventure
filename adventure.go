package adventure

import (
	"fmt"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

func Run() {
	var (
		err    error
		client *twittergo.Client
		resp   *twittergo.APIResponse
	)

	// set up game world
	// set up twitter client for adventure user
	client, err = loadCredentials()
	if err != nil {
		log.Fatal(err)
	}

	// print some debug on the user
	user := &twittergo.User{}
	resp = doRequest(client, "/1.1/account/verify_credentials.json")
	parseWithErrorHandling(resp, user)

	fmt.Printf("ID:                   %v\n", user.Id())
	fmt.Printf("Name:                 %v\n", user.Name())
	printResponseRateLimits(resp)

	// setup channel for listen loop to tell game loop
	// about incoming tweets
	tweetChannel := make(chan *twittergo.Tweet)
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	// setup listen loop for @mentions
	go func() {
		// each tweet mentioned stuff onto channel
		timeline := &twittergo.Timeline{}
		resp = doRequest(client, "/1.1/statuses/mentions_timeline.json")
		parseWithErrorHandling(resp, timeline)
		fmt.Printf("Num Mentions:   %v\n", len(*timeline))
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
		}
		// set gamestate
		// tweet at them their "room"
	}()

	waitGroup.Wait()

}

func doRequest(client *twittergo.Client, api_path string) (resp *twittergo.APIResponse) {
	var (
		req *http.Request
		err error
	)
	req, err = http.NewRequest("GET", api_path, nil)
	if err != nil {
		log.Fatalf("Could not parse request: %v\n", err)
	}
	resp, err = client.SendRequest(req)
	if err != nil {
		log.Fatalf("Could not send request: %v\n", err)
	}
	return
}

func parseWithErrorHandling(resp *twittergo.APIResponse, out interface{}) {
	err := resp.Parse(out)
	if err != nil {
		log.Fatalf("Problem parsing response: %v\n", err)
	}
}

func printResponseRateLimits(resp *twittergo.APIResponse) {
	if resp.HasRateLimit() {
		fmt.Printf("Rate limit:           %v\n", resp.RateLimit())
		fmt.Printf("Rate limit remaining: %v\n", resp.RateLimitRemaining())
		fmt.Printf("Rate limit reset:     %v\n", resp.RateLimitReset())
	} else {
		fmt.Printf("Could not parse rate limit from response.\n")
	}
}

func loadCredentials() (client *twittergo.Client, err error) {
	credentials, err := ioutil.ReadFile("CREDENTIALS")
	if err != nil {
		log.Fatal("CREDENTIALS file missing")
	}
	lines := strings.Split(string(credentials), "\n")
	config := &oauth1a.ClientConfig{
		ConsumerKey:    lines[0],
		ConsumerSecret: lines[1],
	}
	user := oauth1a.NewAuthorizedConfig(lines[2], lines[3])
	client = twittergo.NewClient(config, user)
	return
}
