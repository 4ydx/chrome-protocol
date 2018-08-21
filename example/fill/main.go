package main

import (
	"github.com/4ydx/chrome-protocol"
	"github.com/4ydx/chrome-protocol/actions/dom"
	"github.com/4ydx/chrome-protocol/actions/enable"
	"github.com/4ydx/chrome-protocol/actions/input"
	"github.com/4ydx/chrome-protocol/actions/page"
	"log"
	"time"
)

func main() {
	frame := cdp.Start(9222)
	defer cdp.Stop()

	// Enable page and dom events
	if err := enable.Page(frame, time.Second*2); err != nil {
		panic(err)
	}
	if err := enable.Dom(frame, time.Second*2); err != nil {
		panic(err)
	}

	// Navigate
	if err := page.Navigate(frame, "https://google.com", time.Second*5); err != nil {
		panic(err)
	}

	// Focus
	if err := dom.Focus(frame, "#lst-ib", time.Second*5); err != nil {
		panic(err)
	}

	// Fill
	if err := input.Fill(frame, "testing", time.Second*5); err != nil {
		panic(err)
	}

	log.Printf("\n-- All completed for %s --\n", frame.FrameID)
}
