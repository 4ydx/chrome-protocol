package actions

import (
	"github.com/4ydx/cdp/protocol/dom"
	"github.com/4ydx/cdp/protocol/page"
	"github.com/4ydx/chrome-protocol"
	"strings"
	"testing"
	"time"
)

func TestClick(t *testing.T) {
	browser := cdp.NewBrowser(BrowserPath, 9222, "dom_test.log")

	frame := cdp.Start(browser, cdp.LogBasic)
	defer frame.Stop(true)

	// Enable page, dom, and network events
	if err := EnablePage(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
	if err := EnableDom(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
	if err := EnableNetwork(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}

	// Navigate
	if _, err := Navigate(frame, "https://google.com", time.Second*5); err != nil {
		t.Fatal(err)
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
		cdp.Event{Name: page.EventPageNavigatedWithinDocument, Value: &page.NavigatedWithinDocumentReply{}, IsRequired: true},
	}
	events = append(events, GetNavigationEvents()...)
	events, err := Click(frame, "gb_70", events, time.Second*5)
	if err != nil {
		t.Fatal(err)
	}
	url := events[0].Value.(*page.NavigatedWithinDocumentReply).URL
	if !strings.HasPrefix(url, "https://accounts.google.com/signin") {
		t.Fatalf("Missing prefix 'https://accounts.google.com/signin' in resulting url %s", url)
	}
	t.Logf("All completed for %s", frame.FrameID)
}

func TestDOMClearedWhenEventSpecified(t *testing.T) {
	browser := cdp.NewBrowser(BrowserPath, 9222, "dom_test.log")

	frame := cdp.Start(browser, cdp.LogAll)
	defer frame.Stop(true)

	// Enable page, dom, and network events
	if err := EnablePage(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
	if err := EnableDom(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
	if err := EnableNetwork(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}

	// Navigate
	if _, err := Navigate(frame, "https://google.com", time.Second*5); err != nil {
		t.Fatal(err)
	}

	// Click triggers navitation to another page.
	events := []cdp.Event{
		cdp.Event{Name: page.EventPageNavigatedWithinDocument, Value: &page.NavigatedWithinDocumentReply{}, IsRequired: true},

		// Specifying the DocumentUpdated event.
		cdp.Event{Name: dom.EventDOMDocumentUpdated, Value: &dom.DocumentUpdatedReply{}, IsRequired: true},
	}
	events = append(events, GetNavigationEvents()...)

	events, err := Click(frame, "gb_70", events, time.Second*5)
	if err != nil {
		t.Fatal(err)
	}
	url := events[0].Value.(*page.NavigatedWithinDocumentReply).URL
	if !strings.HasPrefix(url, "https://accounts.google.com/signin") {
		t.Fatalf("Missing prefix 'https://accounts.google.com/signin' in resulting url %s", url)
	}
	t.Logf("All completed for %s", frame.FrameID)

	// Check to see that the DOM has been emptied.  The new page load should trigger a dom.documentUpdated event.
	if frame.GetDOM() != nil {
		t.Fatal("Expecting the Frame.DOM object to be set to nil due to a dom.documentUpdated event.")
	}
}

func TestDOMClearedWhenEventNotSpecified(t *testing.T) {
	browser := cdp.NewBrowser(BrowserPath, 9222, "dom_test.log")

	frame := cdp.Start(browser, cdp.LogAll)
	defer frame.Stop(true)

	// Enable page, dom, and network events
	if err := EnablePage(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
	if err := EnableDom(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
	if err := EnableNetwork(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}

	// Navigate
	if _, err := Navigate(frame, "https://google.com", time.Second*5); err != nil {
		t.Fatal(err)
	}

	// Click triggers navitation to another page.
	events := []cdp.Event{
		cdp.Event{Name: page.EventPageNavigatedWithinDocument, Value: &page.NavigatedWithinDocumentReply{}, IsRequired: true},
	}
	events = append(events, GetNavigationEvents()...)

	events, err := Click(frame, "gb_70", events, time.Second*5)
	if err != nil {
		t.Fatal(err)
	}
	url := events[0].Value.(*page.NavigatedWithinDocumentReply).URL
	if !strings.HasPrefix(url, "https://accounts.google.com/signin") {
		t.Fatalf("Missing prefix 'https://accounts.google.com/signin' in resulting url %s", url)
	}
	t.Logf("All completed for %s", frame.FrameID)

	// Check to see that the DOM has been emptied.  The new page load should trigger a dom.documentUpdated event.
	if frame.GetDOM() != nil {
		t.Fatal("Expecting the Frame.DOM object to be set to nil due to a dom.documentUpdated event.")
	}
}
