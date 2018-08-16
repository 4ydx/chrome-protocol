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
	id := cdp.Start()

	// Enable dom and page communication with chrome
	ea.EnableDom(id, time.Second*2).Run()
	ea.EnablePage(id, time.Second*2).Run()

	// Navigate
	a0 := pa.Navigate(id, "https://google.com", time.Second*5)
	a0.Run()

	// FindAll objects matching the given string
	/*
		_, res0 := da.FindAll(id, "#lst-ib", time.Second*5)
		if len(res0.NodeIds) == 0 {
			panic("Expecting the search field to be present.")
		}
	*/
	res1 := da.Focus(id, "#lst-ib", time.Second*5)

	log.Print("\n-- All completed --\n")
	a0.Log()
	//log.Printf("res0 %+v\n", res0)
	log.Printf("res1 %+v\n", res1)

	cdp.Stop()
}
