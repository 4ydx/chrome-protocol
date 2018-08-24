package main

import (
	"fmt"
	ppage "github.com/4ydx/cdp/protocol/page"
	"github.com/4ydx/chrome-protocol"
	"github.com/4ydx/chrome-protocol/actions/dom"
	"github.com/4ydx/chrome-protocol/actions/enable"
	"github.com/4ydx/chrome-protocol/actions/page"
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
	//
	// In addition, there is a "Page.navigatedWithinDocument" which is the page and url that is ultimately displayed.
	// The only way to be able to see this is to run this binary (./click) and then look at the resulting "log.txt" file.
	// Internally you will find the event firing while the login page frame is being loaded.
	events := []cdp.Event{
		cdp.Event{Name: ppage.EventPageNavigatedWithinDocument, Value: &ppage.NavigatedWithinDocumentReply{}, IsRequired: true},
	}
	events = append(events, page.GetNavigationEvents()...)
	events, err := dom.Click(frame, "gb_70", events, time.Second*5)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", events[0].Value.(*ppage.NavigatedWithinDocumentReply).URL)

	log.Printf("\n-- All completed for %s --\n", frame.FrameID)
}
