package main

import (
	"github.com/4ydx/chrome-protocol"
	"github.com/4ydx/chrome-protocol/actions/dom"
	"github.com/4ydx/chrome-protocol/actions/enable"
	"github.com/4ydx/chrome-protocol/actions/page"
	"log"
	"time"
)

func main() {
	eventCache, id, actionChan, stepChan, allComplete, shutdown := cdp.Start()

	// Enable all communication with chrome
	a0 := ea.EnableDom(id, time.Second*2)
	a0.Run(eventCache, actionChan, stepChan)
	a1 := ea.EnablePage(id, time.Second*2)
	a1.Run(eventCache, actionChan, stepChan)

	// Navigate
	a2 := pa.Navigate(id, "https://google.com", time.Second*5)
	a2.Run(eventCache, actionChan, stepChan)

	// FindAll objects matching the given string
	res0 := da.FindAll(id, "lst-ib", time.Second*5, eventCache, actionChan, stepChan)
	res1 := da.FindAll(id, "hplogo", time.Second*5, eventCache, actionChan, stepChan)

	log.Print("\n-- All completed --\n")
	a0.Log()
	a1.Log()
	a2.Log()
	log.Printf("res0 %+v\n", res0)
	log.Printf("res1 %+v\n", res1)

	time.Sleep(time.Second * 10)
	cdp.Stop(allComplete, shutdown)
}
