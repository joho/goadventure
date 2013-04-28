package main

import (
	"flag"
	"github.com/joho/goadventure"
	"log"
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

	log.Println("Server starting up. SIGINT (CTRL+C) to quit.")
	// Use all the threads.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Setup go routine for signal management
	stopRunning := make(chan bool)
	go func() {
		signalChannel := make(chan os.Signal)
		signal.Notify(signalChannel, syscall.SIGINT)

		<-signalChannel
		log.Println("\n\nSIGINT Received, Shutting Down")
		stopRunning <- true
		close(stopRunning)
	}()

	// Command line flags
	useLiveTwitterClient := flag.Bool("live-twitter", false, "set to actually talk to live twitter")
	usePersistentStorage := flag.Bool("persistent-storage", false, "set to use persistent storage for game state")
	flag.Parse()

	if *useLiveTwitterClient {
		log.Println("Using Twitter API for input/output")
		twitterWrapper = goadventure.NewRealTwitterWrapper()
	} else {
		log.Println("Using interactive input/output")
		twitterWrapper = goadventure.NewFakeTwitterWrapper()
	}

	if *usePersistentStorage {
		log.Println("Using persistent storage for game state")
		storageEngine = goadventure.NewPersistentStorageEngine()
	} else {
		log.Println("Using in memory storage for game state")
		storageEngine = goadventure.NewInMemoryStorageEngine()
	}

	// Let's play!
	goadventure.Run(stopRunning, twitterWrapper, storageEngine)
}
