package goadventure

import (
	"fmt"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type TwitterWrapper struct {
	client *twittergo.Client
}

func NewTwitterWrapper() *TwitterWrapper {
	client := loadCredentials()
	twitterWrapper := TwitterWrapper{client}
	return &twitterWrapper
}

func (twitterWrapper TwitterWrapper) PrintUserDebugInfo() {
	var resp *twittergo.APIResponse
	user := &twittergo.User{}
	resp = twitterWrapper.doRequest("/1.1/account/verify_credentials.json")
	parseWithErrorHandling(resp, user)

	fmt.Printf("ID:                   %v\n", user.Id())
	fmt.Printf("Name:                 %v\n", user.Name())
	printResponseRateLimits(resp)
}

func (twitterWrapper TwitterWrapper) GetUserMentionsTimeline() (timeline *twittergo.Timeline) {
	var resp *twittergo.APIResponse
	timeline = &twittergo.Timeline{}
	resp = twitterWrapper.doRequest("/1.1/statuses/mentions_timeline.json")
	parseWithErrorHandling(resp, timeline)
	fmt.Printf("Num Mentions:   %v\n", len(*timeline))
	return
}

func (twitterWrapper TwitterWrapper) SendResponseToUser(user *twittergo.User, message string) {
	fmt.Printf("Hypothetically sending '%v' to '%v'", message, user.ScreenName())
}

func (twitterWrapper TwitterWrapper) doRequest(api_path string) (resp *twittergo.APIResponse) {
	var (
		req *http.Request
		err error
	)
	req, err = http.NewRequest("GET", api_path, nil)
	if err != nil {
		log.Fatalf("Could not parse request: %v\n", err)
	}
	resp, err = twitterWrapper.client.SendRequest(req)
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

func loadCredentials() *twittergo.Client {
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
	return twittergo.NewClient(config, user)
}
