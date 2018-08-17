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

	// Enable communication with chrome
	if err := enable.Page(id, time.Second*2); err != nil {
		panic(err)
	}
	if err := enable.Dom(id, time.Second*2); err != nil {
		panic(err)
	}
	if err := enable.Network(id, time.Second*2); err != nil {
		panic(err)
	}

	// Navigate
	if err := page.Navigate(id, "https://google.com", time.Second*5); err != nil {
		panic(err)
	}

	// Click
	if err := dom.Click(id, "gb_70", time.Second*5); err != nil {
		panic(err)
	}

	log.Print("\n-- All completed --\n")

	time.Sleep(10 * time.Second)

	cdp.Stop()
}
