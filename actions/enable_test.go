package actions

import (
	"github.com/4ydx/chrome-protocol"
	"testing"
	"time"
)

func TestEnable(t *testing.T) {
	browser := cdp.NewBrowser("/usr/bin/google-chrome", 9222)

	frame := cdp.Start(9222, cdp.LogBasic)
	defer func() {
		cdp.Stop()
		browser.Stop()
	}()

	if err := EnablePage(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
	if err := EnableDom(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
	if err := EnableRuntime(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
	if err := EnableNetwork(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
	if err := EnableCSS(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
	if err := EnableIndexedDB(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
	if err := EnableAll(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
}
