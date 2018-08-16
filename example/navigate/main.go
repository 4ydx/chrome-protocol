package main

import (
	"github.com/4ydx/chrome-protocol"
	"github.com/4ydx/chrome-protocol/actions/enable"
	"github.com/4ydx/chrome-protocol/actions/page"
	"log"
	"time"
)

func main() {
	eventCache, id, actionChan, stepChan, allComplete, shutdown := cdp.Start()

	// Enable all communication with chrome
	a0 := ea.EnablePage(id, time.Second*2)
	a0.Run(eventCache, actionChan, stepChan)

	// Navigate
	a1 := pa.Navigate(id, "https://google.com", time.Second*5)
	a1.Run(eventCache, actionChan, stepChan)

	log.Print("\n-- All completed --\n")
	a0.Log()
	a1.Log()

	cdp.Stop(allComplete, shutdown)
}
