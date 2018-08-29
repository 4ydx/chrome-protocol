[![](https://godoc.org/github.com/4ydx/chrome-protocol?status.svg)](http://godoc.org/github.com/4ydx/chrome-protocol)

# About chrome-protocol

A relatively thin wrapper on top of code that is generated based on
the chrome devtool protocol.  Aims to provide a few of the basic commands that
one would desire when automating actions in chrome or any other browser that
supports the protocol.

This is still a work in progress.

- Very fast.
- No hidden errors.
- Simple approach makes it easy to understand what is happening under the hood.

## Examples

Look under [github.com/4ydx/chrome-protocol/actions](https://github.com/4ydx/chrome-protocol/tree/master/actions).  The testing files are the examples.  There is one example in the example folder.

- Click
- Fill
- Focus
- Navigation
- Screenshot
- As well as other actions (css style retrieval, javascript evaluation, etc).

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
	browser := cdp.NewBrowser("/usr/bin/google-chrome", 9222, "browser.log")

	frame := cdp.Start(browser, cdp.LogBasic)
	defer func() {
		// passing false prevents the browser from stopping immediately
		frame.Stop(false)

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

Actions encapsulate everything you need in order to interact with a browser. An action contains commands and events.

When you construct an action, you need to fill in at least one command that consists of the struct representing the parameters that will be sent with the command,
the struct that represents the reply to that command from the server, and the method name of the API call you are making.

It is possible to associate events that the server will send to the client with your action.  By specifying events you can be sure that a given action has actually run its
course and the browser state is where you would expect it to be.

API methods, command parameters, command responses, possible events, and types are all defined in the [Devtools Reference](https://chromedevtools.github.io/devtools-protocol/tot).

This is a possible Navigation action that watches for the FrameStoppedLoadingReply event which helps to ensure that navigation to a url is fully completed.

```
// Navigate sends the browser to the given URL
func Navigate(frame *cdp.Frame, url string, timeout time.Duration) ([]cdp.Event, error) {
	events := GetNavigationEvents()
	action := cdp.NewAction(frame,
		events,
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: page.CommandPageNavigate, Params: &page.NavigateArgs{URL: url}, Reply: &page.NavigateReply{}, Timeout: timeout},
		})
	if err := action.Run(); err != nil {
		log.Print(err)
		return events, err
	}
	return events, nil
}
```

## Caveats

- Concurrent actions are currently not supported.
