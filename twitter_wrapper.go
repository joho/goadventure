package goadventure

import (
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"io/ioutil"
	"log"
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
