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
	if err := ea.EnableDom(id, time.Second*2); err != nil {
		panic(err)
	}
	if err := ea.EnablePage(id, time.Second*2); err != nil {
		panic(err)
	}

	// Navigate
	if err := pa.Navigate(id, "https://google.com", time.Second*5); err != nil {
		panic(err)
	}

	// Focus
	if err := da.Focus(id, "#lst-ib", time.Second*5); err != nil {
		panic(err)
	}

	log.Print("\n-- All completed --\n")
	cdp.Stop()
}
