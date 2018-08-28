package actions

import (
	"github.com/4ydx/chrome-protocol"
	"testing"
	"time"
)

func TestFill(t *testing.T) {
	browser := cdp.NewBrowser(BrowserPath, 9222, "input_test.log")

	frame := cdp.Start(browser, cdp.LogBasic)
	defer frame.Stop(true)

	// Enable page and dom events
	if err := EnablePage(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
	if err := EnableDom(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}

	// Navigate
	if _, err := Navigate(frame, "https://google.com", time.Second*5); err != nil {
		t.Fatal(err)
	}

	// Fill
	if err := Fill(frame, "#lst-ib", "testing", time.Second*5); err != nil {
		t.Fatal(err)
	}

	// Run Script
	script := "document.getElementById('lst-ib').value"
	reply, err := Evaluate(frame, script, time.Second*5)
	if err != nil {
		t.Fatal(err)
	}
	if string(*reply.Result.Value) != "\"testing\"" {
		t.Fatalf("Expecting \"testing\" but got %s", *reply.Result.Value)
	}
	t.Logf("All completed for %s", frame.FrameID)
}
