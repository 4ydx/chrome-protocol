package main

import (
	"github.com/4ydx/chrome-protocol"
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

	// Enable page events
	if err := enable.Page(frame, time.Second*2); err != nil {
		panic(err)
	}

	// Navigate
	if err := page.Navigate(frame, "https://google.com", time.Second*5); err != nil {
		panic(err)
	}

	// Screenshot
	if err := page.Screenshot(frame, "file", "jpeg", 100, nil, time.Second*5); err != nil {
		panic(err)
	}

	log.Printf("\n-- All completed for %s --\n", frame.FrameID)
}
