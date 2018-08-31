package actions

import (
	"github.com/4ydx/chrome-protocol"
	"testing"
	"time"
)

func TestCookies(t *testing.T) {
	browser := cdp.NewBrowser(BrowserPath, 9222, "page_test.log")

	frame := cdp.Start(browser, cdp.LogBasic)
	defer frame.Stop(true)

	// Enable page events
	if err := EnablePage(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
	if err := EnableNetwork(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}

	// Navigate
	_, err := Navigate(frame, "https://google.com", time.Second*5)
	if err != nil {
		t.Fatal(err)
	}

	// Get Cookies
	cookies, err := Cookies(frame, time.Second*5)
	if err != nil {
		t.Fatal(err)
	}
	if len(cookies) == 0 {
		t.Fatal("No cookies found.")
	}
	t.Logf("All completed for %s", frame.FrameID)
}
