package actions

import (
	"github.com/4ydx/chrome-protocol"
	"testing"
	"time"
)

func TestEvaluate(t *testing.T) {
	browser := cdp.NewBrowser("/usr/bin/google-chrome", 9222)

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

	// Run Script
	script := "document.getElementsByName('btnI')[0].value"
	reply, err := Evaluate(frame, script, time.Second*5)
	if err != nil {
		t.Fatal(err)
	}
	if string(*reply.Result.Value) != "\"I'm Feeling Lucky\"" {
		t.Fatalf("Expecting 'I'm Feeling Lucky' but got %s", *reply.Result.Value)
	}
	t.Logf("All completed for %s", frame.FrameID)
}
