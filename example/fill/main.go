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
	id := cdp.Start()

	// Enable all communication with chrome
	if err := enable.Page(id, time.Second*2); err != nil {
		panic(err)
	}
	if err := enable.Dom(id, time.Second*2); err != nil {
		panic(err)
	}

	// Navigate
	if err := page.Navigate(id, "https://google.com", time.Second*5); err != nil {
		panic(err)
	}

	// Focus
	if err := dom.Focus(id, "#lst-ib", time.Second*5); err != nil {
		panic(err)
	}

	// Fill
	if err := input.Fill(id, "testing", time.Second*5); err != nil {
		panic(err)
	}

	log.Print("\n-- All completed --\n")

	cdp.Stop()
}
