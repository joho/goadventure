package main

import (
	"fmt"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	var (
		err           error
		client        *twittergo.Client
		resp          *twittergo.APIResponse
	)

	client, err = LoadCredentials()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

  user := &twittergo.User{}
	resp = DoRequest(client, "/1.1/account/verify_credentials.json")
	ParseWithErrorHandling(resp, user)

	fmt.Printf("ID:                   %v\n", user.Id())
	fmt.Printf("Name:                 %v\n", user.Name())
	PrintResponseRateLimits(resp)

  timeline := &twittergo.Timeline{}
	resp = DoRequest(client, "/1.1/statuses/mentions_timeline.json")
	ParseWithErrorHandling(resp, timeline)
	fmt.Printf("Num Mentions:   %v\n", len(*timeline))
	for _, tweet := range *timeline {
    fmt.Printf("Tweet:   %v\n", tweet.Text())
	}

}

func DoRequest(client *twittergo.Client, api_path string) (resp *twittergo.APIResponse) {
	var (
		req *http.Request
		err error
	)
	req, err = http.NewRequest("GET", api_path, nil)
	if err != nil {
		fmt.Printf("Could not parse request: %v\n", err)
		os.Exit(1)
	}
	resp, err = client.SendRequest(req)
	if err != nil {
		fmt.Printf("Could not send request: %v\n", err)
		os.Exit(1)
	}
	return
}

func ParseWithErrorHandling(resp *twittergo.APIResponse, out interface{}) {
	err := resp.Parse(out)
	if err != nil {
		fmt.Printf("Problem parsing response: %v\n", err)
		os.Exit(1)
	}
}

func PrintResponseRateLimits(resp *twittergo.APIResponse) {
	if resp.HasRateLimit() {
		fmt.Printf("Rate limit:           %v\n", resp.RateLimit())
		fmt.Printf("Rate limit remaining: %v\n", resp.RateLimitRemaining())
		fmt.Printf("Rate limit reset:     %v\n", resp.RateLimitReset())
	} else {
		fmt.Printf("Could not parse rate limit from response.\n")
	}
}

func LoadCredentials() (client *twittergo.Client, err error) {
	credentials, err := ioutil.ReadFile("CREDENTIALS")
	if err != nil {
		return
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
