package main

import (
	"github.com/4ydx/chrome-protocol"
	"github.com/4ydx/chrome-protocol/actions"
	"log"
	"time"
)

func main() {
	browser := cdp.NewBrowser("/usr/bin/google-chrome", 9222)

	frame := cdp.Start(9222, cdp.LOG_BASIC)
	defer func() {
		cdp.Stop()

		// Give yourself time to view the final page in the browser.
		time.Sleep(3 * time.Second)
		browser.Stop()
	}()

	// Enable page events
	if err := actions.EnablePage(frame, time.Second*2); err != nil {
		panic(err)
	}

	// Navigate
	if _, err := actions.Navigate(frame, "https://google.com", time.Second*5); err != nil {
		panic(err)
	}

	// Screenshot
	if err := actions.Screenshot(frame, "screenshot", "jpeg", 100, nil, time.Second*5); err != nil {
		panic(err)
	}

	log.Printf("\n-- All completed for %s --\n", frame.FrameID)
}
