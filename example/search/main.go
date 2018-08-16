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
	actionCache, id, actionChan, stepChan, allComplete, shutdown := cdp.Start()

	// Enable all communication with chrome
	a0 := ea.EnableDom(id, time.Second*2)
	a0.Run(actionCache, actionChan, stepChan)
	a1 := ea.EnablePage(id, time.Second*2)
	a1.Run(actionCache, actionChan, stepChan)

	// Navigate
	a2 := pa.Navigate(id, "https://google.com", time.Second*5)
	a2.Run(actionCache, actionChan, stepChan)

	// FindAll objects matching the given string
	res0 := da.FindAll(id, "#lst-ib", time.Second*5, actionCache, actionChan, stepChan)
	if len(res0.NodeIds) == 0 {
		panic("Expecting the search field to be present.")
	}
	res1 := da.Focus(id, res0.NodeIds[0], time.Second*5, actionCache, actionChan, stepChan)

	log.Print("\n-- All completed --\n")
	a0.Log()
	a1.Log()
	a2.Log()
	log.Printf("res0 %+v\n", res0)
	log.Printf("res1 %+v\n", res1)

	time.Sleep(time.Second * 10)
	cdp.Stop(allComplete, shutdown)
}
