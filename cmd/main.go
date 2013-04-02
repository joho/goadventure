package main

import (
	"flag"
	"fmt"
	"github.com/joho/goadventure"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	useLiveTwitterClient := flag.Bool("live-twitter", false, "set to actually talk to live twitter")
	flag.Parse()

	fmt.Println("Server starting up. CTRL+C to quit (may take a few seconds)")
	stopRunning := make(chan bool)
	go func() {
		signalChannel := make(chan os.Signal)
		signal.Notify(signalChannel, syscall.SIGINT)

		<-signalChannel
		fmt.Println("\n\nCTRL+C Received, will eventually halt")
		stopRunning <- true
	}()

	if *useLiveTwitterClient {
		fmt.Println("Using real twitter wrapper (-live-twitter has been set)")
		goadventure.Run(stopRunning, goadventure.NewRealTwitterWrapper())
	} else {
		fmt.Println("Using fake twitter wrapper (use -live-twitter flag to do fo' reals)")
		goadventure.Run(stopRunning, new(goadventure.FakeTwitterWrapper))
	}
}
