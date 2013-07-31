package goadventure

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type RealTwitterWrapper struct {
	client *twittergo.Client
}

func NewRealTwitterWrapper() *RealTwitterWrapper {
	client := loadCredentials()
	tw := RealTwitterWrapper{client}
	tw.printUserDebugInfo()
	return &tw
}

func (tw *RealTwitterWrapper) DurationUntilNextRead() time.Duration {
	// TODO make this dynamically look at rate limit stuff
	return 1 * time.Minute
}

func (tw *RealTwitterWrapper) printUserDebugInfo() {
	var resp *twittergo.APIResponse
	user := &twittergo.User{}
	resp = tw.doGetRequest("/1.1/account/verify_credentials.json")
	parseWithErrorHandling(resp, user)

	log.Printf("ID:                   %v\n", user.Id())
	log.Printf("Name:                 %v\n", user.Name())
	printResponseRateLimits(resp)
}

func (tw *RealTwitterWrapper) GetUserMentionsTimeline(tweetChannel chan *twittergo.Tweet) {
	var (
		resp     *twittergo.APIResponse
		timeline *twittergo.Timeline
	)

	timeline = &twittergo.Timeline{}

	resp = tw.doGetRequest("/1.1/statuses/mentions_timeline.json")
	parseWithErrorHandling(resp, timeline)

	log.Printf("Num Mentions:   %v\n", len(*timeline))

	for _, tweet := range *timeline {
		tweetChannel <- &tweet
	}
}

func (tw *RealTwitterWrapper) RespondToTweet(tweet *twittergo.Tweet, message string) {
	var (
		err  error
		user twittergo.User
		req  *http.Request
		resp *twittergo.APIResponse
	)

	user = tweet.User()
	data := url.Values{}

	// set status
	status := fmt.Sprintf("@%v %v", user.ScreenName(), message)
	data.Set("status", status)
	// set in_reply_to_status_id
	status_id := fmt.Sprintf("%d", tweet.Id())
	data.Set("in_reply_to_status_id", status_id)

	log.Printf("Set status '%v' to '%v' in reply to %v", status, user.ScreenName(), status_id)

	body := strings.NewReader(data.Encode())
	req, err = http.NewRequest("POST", "/1.1/statuses/update.json", body)
	if err != nil {
		log.Fatalf("Could not parse request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = tw.client.SendRequest(req)
	if err != nil {
		log.Fatalf("Could not send request: %v\n", err)
	}
	tweet = &twittergo.Tweet{}
	err = resp.Parse(tweet)
	if err != nil {
		log.Fatalf("Problem parsing response: %v\n", err)
	}
}

func (tw *RealTwitterWrapper) doGetRequest(api_path string) (resp *twittergo.APIResponse) {
	var (
		req *http.Request
		err error
	)
	req, err = http.NewRequest("GET", api_path, nil)
	if err != nil {
		log.Fatalf("Could not parse request: %v\n", err)
	}
	resp, err = tw.client.SendRequest(req)
	if err != nil {
		log.Fatalf("Could not send request: %v\n", err)
	}
	return
}

func parseWithErrorHandling(resp *twittergo.APIResponse, out interface{}) {
	err := resp.Parse(out)
	// TODO recover from 2013/04/28 14:16:32 Problem parsing response: Error 187: Status is a duplicate.
	if err != nil {
		log.Fatalf("Problem parsing response: %v\n", err)
	}
}

func printResponseRateLimits(resp *twittergo.APIResponse) {
	if resp.HasRateLimit() {
		log.Printf("Rate limit:           %v\n", resp.RateLimit())
		log.Printf("Rate limit remaining: %v\n", resp.RateLimitRemaining())
		log.Printf("Rate limit reset:     %v\n", resp.RateLimitReset())
	} else {
		log.Printf("Could not parse rate limit from response.\n")
	}
}

func loadCredentials() *twittergo.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file missing")
	}

	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	config := &oauth1a.ClientConfig{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
	}

	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	user := oauth1a.NewAuthorizedConfig(accessToken, accessTokenSecret)

	return twittergo.NewClient(config, user)
}
