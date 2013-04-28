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
	var (
		twitterWrapper goadventure.TwitterWrapper
		storageEngine  goadventure.StorageEngine
	)

	fmt.Println("Server starting up. SIGINT (CTRL+C) to quit.")
	// Use all the threads.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Setup go routine for signal management
	stopRunning := make(chan bool)
	go func() {
		signalChannel := make(chan os.Signal)
		signal.Notify(signalChannel, syscall.SIGINT)

		<-signalChannel
		fmt.Println("\n\nSIGINT Received, Shutting Down")
		stopRunning <- true
		close(stopRunning)
	}()

	// Command line flags
	useLiveTwitterClient := flag.Bool("live-twitter", false, "set to actually talk to live twitter")
	usePersistentStorage := flag.Bool("persistent-storage", false, "set to use persistent storage for game state")
	flag.Parse()

	if *useLiveTwitterClient {
		fmt.Println("Using Twitter API for input/output")
		twitterWrapper = goadventure.NewRealTwitterWrapper()
	} else {
		fmt.Println("Using interactive input/output")
		twitterWrapper = new(goadventure.FakeTwitterWrapper)
	}

	if *usePersistentStorage {
		fmt.Println("Using persistent storage for game state")
		storageEngine = goadventure.NewPersistentStorageEngine()
	} else {
		fmt.Println("Using in memory storage for game state")
		storageEngine = goadventure.NewInMemoryStorageEngine()
	}

	// Let's play!
	goadventure.Run(stopRunning, twitterWrapper, storageEngine)
}
