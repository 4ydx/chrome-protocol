package main

import (
	"github.com/4ydx/chrome-protocol"
	"github.com/4ydx/chrome-protocol/actions/enable"
	"github.com/4ydx/chrome-protocol/actions/page"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	f, err := os.Create("log.txt")
	if err != nil {
		panic(err)
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetOutput(f)

	c := cdp.GetWebsocket()
	defer c.Close()

	actionChan := make(chan *cdp.Action)
	stepChan := make(chan bool)

	shutdown := make(chan struct{})
	allComplete := make(chan struct{})

	stepCache := cdp.NewStepCache()
	eventCache := cdp.NewEventCache()

	go cdp.Write(c, actionChan, stepCache, shutdown, allComplete)
	go cdp.Read(c, stepChan, stepCache, eventCache, shutdown)

	id := &cdp.ID{
		RWMutex: &sync.RWMutex{},
		Value:   11111,
	}

	// Enable all communication with chrome
	a0 := ea.EnablePage(id, time.Second*2)
	a0.Run(eventCache, actionChan, stepChan)

	// Navigate
	a1 := pa.Navigate(id, "https://google.com", time.Second*5)
	a1.Run(eventCache, actionChan, stepChan)

	log.Print("\n-- All completed --\n")
	a0.Log()
	a1.Log()

	allComplete <- struct{}{}
	<-shutdown
}
