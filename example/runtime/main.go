package main

import (
	"fmt"
	"github.com/4ydx/chrome-protocol"
	"github.com/4ydx/chrome-protocol/actions"
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

	// Enable page and dom events
	if err := actions.EnablePage(frame, time.Second*2); err != nil {
		panic(err)
	}
	if err := actions.EnableDom(frame, time.Second*2); err != nil {
		panic(err)
	}
	if err := actions.EnableRuntime(frame, time.Second*2); err != nil {
		panic(err)
	}

	// Navigate
	if err := actions.Navigate(frame, "https://google.com", time.Second*5); err != nil {
		panic(err)
	}

	// Fill
	if err := actions.Fill(frame, "#lst-ib", "testing", time.Second*5); err != nil {
		panic(err)
	}

	// Determine if the input has the value "testing"
	reply, err := actions.Evaluate(frame, "document.getElementById('lst-ib').value.toString()", time.Second*5)
	if err != nil {
		panic(err)
	}

	// We expect that the value of the element will be "testing".  I am seeing a value of &"testing" and I'm not familiar
	// enough to know why this is prepended by an ampersand, but it is there.  The actual data in the byte array is only
	// the values "testing".
	fmt.Printf("%s\n", reply.Result.Value)

	log.Printf("\n-- All completed for %s --\n", frame.FrameID)
}
