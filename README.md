[![](https://godoc.org/github.com/4ydx/chrome-protocol?status.svg)](http://godoc.org/github.com/4ydx/chrome-protocol)

# About chrome-protocol

A relatively thin wrapper on top of code that is generated based on
the chrome devtool protocol.  Aims to provide a few of the basic commands that
one would desire when automating actions in chrome or any other browser that
supports the protocol.

This is still a work in progress.

## Examples

Look under [github.com/4ydx/chrome-protocol/actions](https://github.com/4ydx/chrome-protocol/actions).  The testing files are the examples.

- Click
- Fill
- Focus
- Navigation
- Screenshot

I will be working on other actions as I need them for my own personal projects.  

You can take the generated code in [github.com/4ydx/cdp](https://github.com/4ydx/cdp/tree/master/protocol) and create your own higher level actions for
interacting with the browser.  This will require understanding the [Devtools Reference](https://chromedevtools.github.io/devtools-protocol/tot).

Navigation example:

```
package main

import (
	"github.com/4ydx/chrome-protocol"
	"github.com/4ydx/chrome-protocol/actions"
	"log"
	"time"
)

func main() {
	browser := cdp.NewBrowser("/usr/bin/google-chrome", 9222)

	frame := cdp.Start(9222, cdp.LOG_BASIC)
	defer func() {
		cdp.Stop()

		// Give yourself time to view the final page in the browser.
		time.Sleep(3 * time.Second)
		browser.Stop()
	}()

	// Enable page events
	if err := actions.EnablePage(frame, time.Second*2); err != nil {
		panic(err)
	}

	// Navigate
	if _, err := actions.Navigate(frame, "https://google.com", time.Second*5); err != nil {
		panic(err)
	}

	log.Printf("\n-- All completed for %s --\n", frame.FrameID)
}
```

## Creating your own Actions

Actions are the requests that you make to the browser in order to automate different tasks.  For instance, asking
the browser to navigate to a particular url.  When you construct an action, you need to fill in at least one "step" that consists
of the params struct, the reply struct, and the method name of the API call you are making.  Finally, it is possible to associate events
that the server will send to the client with your action.  By specifying events you can be sure that a given action has actually run its
course and the browser state is where you would expect it to be.

API methods, events, and types are all defined in the [Devtools Reference](https://chromedevtools.github.io/devtools-protocol/tot).

Possible Navigation Method.  This watches for the FrameStoppedLoadingReply event which helps to ensure that navigation is fully completed.

```
func Navigate(frame *cdp.Frame, url string, timeout time.Duration) error {
	return cdp.NewAction(frame,
		[]cdp.Event{
			cdp.Event{Name: page.EventPageFrameNavigated, Value: &page.FrameNavigatedReply{}, IsRequired: true},
			cdp.Event{Name: page.EventPageFrameStartedLoading, Value: &page.FrameStartedLoadingReply{}, IsRequired: true},
			cdp.Event{Name: page.EventPageFrameStoppedLoading, Value: &page.FrameStoppedLoadingReply{}, IsRequired: true},
		},
		[]cdp.Step{
			cdp.Step{ID: frame.RequestID.GetNext(), Method: page.CommandPageNavigate, Params: &page.NavigateArgs{URL: url}, Reply: &page.NavigateReply{}, Timeout: timeout},
		}).Run()
}
```

## Caveats

- Concurrent actions are currently not supported.
