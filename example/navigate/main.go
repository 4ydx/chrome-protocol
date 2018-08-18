package main

import (
	"github.com/4ydx/chrome-protocol"
	"github.com/4ydx/chrome-protocol/actions/enable"
	"github.com/4ydx/chrome-protocol/actions/page"
	"log"
	"time"
)

func main() {
	p := cdp.Start()

	// Enable all communication with chrome
	if err := enable.Page(p, time.Second*2); err != nil {
		panic(err)
	}

	// Navigate
	if err := page.Navigate(p, "https://google.com", time.Second*5); err != nil {
		panic(err)
	}

	log.Print("\n-- All completed --\n")

	cdp.Stop()
}
