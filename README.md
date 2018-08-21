[![](https://godoc.org/github.com/4ydx/chrome-protocol?status.svg)](http://godoc.org/github.com/4ydx/chrome-protocol)

# About chrome-protocol

A relatively thin wrapper on top of code that is generated based on
the chrome devtool protocol.  Aims to provide a few of the basic commands that
one would desire when automating actions in chrome or any other browser that
supports the protocol.

This is still a work in progress.

# Examples

Examples of basic actions are included in the example folder.

- Navigation
- Focus
- Fill
- Click

I will be working on other actions as I need them for my own personal projects.

Navigation example:

```
package main

import (
	"github.com/4ydx/chrome-protocol"
	"github.com/4ydx/chrome-protocol/actions/enable"
	"github.com/4ydx/chrome-protocol/actions/page"
	"log"
	"time"
)

func main() {
	frame := cdp.Start()

	// Enable page events 
	if err := enable.Page(frame, time.Second*2); err != nil {
		panic(err)
	}

	// Navigate
	if err := page.Navigate(frame, "https://google.com", time.Second*5); err != nil {
		panic(err)
	}

	log.Printf("\n-- All completed for %s --\n", frame.FrameID)

	cdp.Stop()
}
```

# Creating your own Actions

Actions are the requests that you make to the browser in order to automate different tasks.  For instance asking
the browser to navigate to a particular url.  When you construct an action you need to fill in at least one "step" that consists
of the params struct, the reply struct, and the method name of the API call you are making.  Finally it is possible to associate events
that the server will send to the client with your action.

API methods, events, and types are all defined in the [Devtools Reference](https://chromedevtools.github.io/devtools-protocol/tot).

Possible Navigation Method:

```
func Navigate(pg *cdp.Frame, url string, timeout time.Duration) error {
	return cdp.NewAction(pg,
		[]cdp.Event{
			cdp.Event{Name: page.EventPageFrameNavigated, Value: &page.FrameNavigatedReply{}, IsRequired: true},
			cdp.Event{Name: page.EventPageFrameStartedLoading, Value: &page.FrameStartedLoadingReply{}, IsRequired: true},
			cdp.Event{Name: page.EventPageFrameStoppedLoading, Value: &page.FrameStoppedLoadingReply{}, IsRequired: true},
		},
		[]cdp.Step{
			cdp.Step{ID: pg.RequestID.GetNext(), Method: page.CommandPageNavigate, Params: &page.NavigateArgs{URL: url}, Reply: &page.NavigateReply{}, Timeout: timeout},
		}).Run()
}
```

# Caveats

- No code yet for managing a browser.
- Concurrent actions are currently not supported.
