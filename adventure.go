package main

import (
  "fmt"
  "log"
  "os"
  "io/ioutil"
  "strings"
  "net/http"
  "github.com/kurrik/oauth1a"
  "github.com/kurrik/twittergo"
)

func main() {
  var (
		err    error
		client *twittergo.Client
		user   *twittergo.User
    resp   *twittergo.APIResponse
	)
  fmt.Println("does this work?")

  client, err = LoadCredentials()
  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }

  resp = DoRequest(client, "/1.1/account/verify_credentials.json")
	user = &twittergo.User{}
	err = resp.Parse(user)
	if err != nil {
		fmt.Printf("Problem parsing response: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("ID:                   %v\n", user.Id())
	fmt.Printf("Name:                 %v\n", user.Name())
	if resp.HasRateLimit() {
		fmt.Printf("Rate limit:           %v\n", resp.RateLimit())
		fmt.Printf("Rate limit remaining: %v\n", resp.RateLimitRemaining())
		fmt.Printf("Rate limit reset:     %v\n", resp.RateLimitReset())
	} else {
		fmt.Printf("Could not parse rate limit from response.\n")
	}
}

func DoRequest(client *twittergo.Client, api_path string) (resp *twittergo.APIResponse) {
  var (
		req    *http.Request
    err   error
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

