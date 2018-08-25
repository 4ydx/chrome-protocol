package main

import (
	"github.com/4ydx/chrome-protocol"
	"github.com/4ydx/chrome-protocol/actions"
	"log"
	"time"
)

func main() {
	browser := cdp.NewBrowser("/usr/bin/google-chrome", 9222)

	frame := cdp.Start(9222)
	defer func() {
		cdp.Stop()

		// Give yourself time to view the final page in the browser.
		time.Sleep(3 * time.Second)
		browser.Stop()
	}()

	// Enable page and dom events
	if err := actions.EnableDom(frame, time.Second*2); err != nil {
		panic(err)
	}
	if err := actions.EnablePage(frame, time.Second*2); err != nil {
		panic(err)
	}

	// Navigate
	if _, err := actions.Navigate(frame, "https://google.com", time.Second*5); err != nil {
		panic(err)
	}

	// Focus
	if err := actions.Focus(frame, "#lst-ib", time.Second*5); err != nil {
		panic(err)
	}

	log.Printf("\n-- All completed for %s --\n", frame.FrameID)
}
