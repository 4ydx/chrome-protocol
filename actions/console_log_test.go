package actions

import (
	"context"
	"github.com/4ydx/chrome-protocol"
	"testing"
	"time"
)

func TestConsoleLog(t *testing.T) {
	srv := LocalServer()

	browser := cdp.NewBrowser(BrowserPath, 9222, "test.log")

	frame := cdp.Start(browser, cdp.LogDetails)
	defer frame.Stop(true)

	if err := EnablePage(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
	if err := EnableRuntime(frame, time.Second*2); err != nil {
		t.Fatal(err)
	}
	if _, err := Navigate(frame, "http://localhost:8080", time.Second*10); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 1)

	if err := srv.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
}
