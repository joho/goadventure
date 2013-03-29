package main

import (
	"flag"
	"github.com/joho/goadventure"
)

func main() {
	useLiveTwitterClient := flag.Bool("live-twitter", false, "set to actually talk to live twitter")
	flag.Parse()

	if *useLiveTwitterClient {
		goadventure.Run(goadventure.NewRealTwitterWrapper())
	} else {
		goadventure.Run(new(goadventure.FakeTwitterWrapper))
	}
}
