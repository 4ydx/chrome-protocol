package actions

import (
	"github.com/4ydx/chrome-protocol"
	"testing"
	"time"
)

func TestEvaluate(t *testing.T) {
	browser := cdp.NewBrowser(BrowserPath, 9222, "evaluate_test.log")

	frame := cdp.Start(browser, cdp.LogBasic)
	defer frame.Stop(true)

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
