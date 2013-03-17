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

type TwitterWrapper interface {
	PrintUserDebugInfo()
	GetUserMentionsTimeline() *twittergo.Timeline
	SendResponseToUser(*twittergo.User, string)
}

type RealTwitterWrapper struct {
	client *twittergo.Client
}

func NewRealTwitterWrapper() *RealTwitterWrapper {
	client := loadCredentials()
	twitterWrapper := RealTwitterWrapper{client}
	return &twitterWrapper
}

func (twitterWrapper *RealTwitterWrapper) PrintUserDebugInfo() {
	var resp *twittergo.APIResponse
	user := &twittergo.User{}
	resp = twitterWrapper.doRequest("/1.1/account/verify_credentials.json")
	parseWithErrorHandling(resp, user)

	fmt.Printf("ID:                   %v\n", user.Id())
	fmt.Printf("Name:                 %v\n", user.Name())
	printResponseRateLimits(resp)
}

func (twitterWrapper *RealTwitterWrapper) GetUserMentionsTimeline() (timeline *twittergo.Timeline) {
	var resp *twittergo.APIResponse
	timeline = &twittergo.Timeline{}
	resp = twitterWrapper.doRequest("/1.1/statuses/mentions_timeline.json")
	parseWithErrorHandling(resp, timeline)
	fmt.Printf("Num Mentions:   %v\n", len(*timeline))
	return
}

func (twitterWrapper *RealTwitterWrapper) SendResponseToUser(user *twittergo.User, message string) {
	fmt.Printf("Hypothetically sending '%v' to '%v'", message, user.ScreenName())
}

func (twitterWrapper *RealTwitterWrapper) doRequest(api_path string) (resp *twittergo.APIResponse) {
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

type FakeTwitterWrapper struct{}

func (twitterWrapper *FakeTwitterWrapper) PrintUserDebugInfo() {
	fmt.Println("I have no actual user, I am pretend")
}

func (twitterWrapper *FakeTwitterWrapper) GetUserMentionsTimeline() *twittergo.Timeline {
	user := map[string]interface{}{
		"screen_name": "johnbarton",
	}
	tweet := twittergo.Tweet{
		"text": "@gotextadventure go north",
		"user": user,
	}
	return &twittergo.Timeline{tweet}
}

func (twitterWrapper *FakeTwitterWrapper) SendResponseToUser(user *twittergo.User, message string) {
	fmt.Printf("Hypothetically Send tweet '%v' to '%v'", message, user.ScreenName())
}
