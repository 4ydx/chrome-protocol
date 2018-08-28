package actions

import (
	"github.com/4ydx/chrome-protocol"
	"testing"
	"time"
)

func TestEnable(t *testing.T) {
	browser := cdp.NewBrowser(BrowserPath, 9222, "enable_test.log")

	frame := cdp.Start(browser, cdp.LogBasic)
	defer frame.Stop(true)

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
