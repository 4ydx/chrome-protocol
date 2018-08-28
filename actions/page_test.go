package actions

import (
	"github.com/4ydx/cdp/protocol/page"
	"github.com/4ydx/chrome-protocol"
	"os"
	"testing"
	"time"
)

func TestNavigate(t *testing.T) {
	browser := cdp.NewBrowser(BrowserPath, 9222)

	frame := cdp.StartWithLog(9222, "page_test.log", cdp.LogBasic)
	defer func() {
		cdp.Stop()
		browser.Stop()
	}()

	// Enable page events
	if err := EnablePage(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}

	// Navigate
	events, err := Navigate(frame, "https://google.com", time.Second*5)
	if err != nil {
		t.Fatal(err)
	}
	url := events[0].Value.(*page.FrameNavigatedReply).Frame.URL
	if url != "https://www.google.com/" {
		t.Fatalf("incorrect url %s", url)
	}
	t.Logf("All completed for %s", frame.FrameID)
}

func TestScreenshot(t *testing.T) {
	browser := cdp.NewBrowser(BrowserPath, 9222)

	frame := cdp.Start(9222, cdp.LogBasic)
	defer func() {
		cdp.Stop()
		browser.Stop()
	}()

	// Enable page events
	if err := EnablePage(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}

	// Navigate
	if _, err := Navigate(frame, "https://google.com", time.Second*5); err != nil {
		t.Fatal(err)
	}

	// Screenshot
	if err := Screenshot(frame, "google", "png", 100, nil, time.Second*5); err != nil {
		t.Fatal(err)
	}
	_, err := os.Stat("google.png")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("All completed for %s", frame.FrameID)
}
