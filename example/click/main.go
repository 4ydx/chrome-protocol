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
	browser := cdp.Browser{}
	browser.Start("/usr/bin/google-chrome", 9222)

	frame := cdp.Start(9222)
	defer func() {
		cdp.Stop()

		// Give yourself time to view the final page in the browser.
		time.Sleep(3 * time.Second)
		browser.Stop()
	}()

	// Enable page, dom, and network events
	if err := enable.Page(frame, time.Second*2); err != nil {
		panic(err)
	}
	if err := enable.Dom(frame, time.Second*2); err != nil {
		panic(err)
	}
	if err := enable.Network(frame, time.Second*2); err != nil {
		panic(err)
	}

	// Navigate
	if err := page.Navigate(frame, "https://google.com", time.Second*5); err != nil {
		panic(err)
	}

	// Click on the google login button which will result in a redirect.
	//
	// Note that we are passing in the required navigation events that will fire as a result of the click.
	// In other words, this click will not be considered completed until the resulting navigation is complete.
	if err := dom.Click(frame, "gb_70", page.GetNavigationEvents(), time.Second*5); err != nil {
		panic(err)
	}

	log.Printf("\n-- All completed for %s --\n", frame.FrameID)
}
