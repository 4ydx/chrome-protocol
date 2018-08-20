[![](https://godoc.org/github.com/4ydx/chrome-protocol?status.svg)](http://godoc.org/github.com/4ydx/chrome-protocol)

# About chrome-protocol

A relatively thin wrapper on top of code that is generated based on
the chrome devtool protocol.  Aims to provide a few of the basic commands that
one would desire when automating actions in chrome or any other browser that
supports the protocol.

# Examples

Examples of basic actions are included in the example folder.

- Navigation
- Focus
- Fill
- Click

I will be working on other actions as I need them for my own personal projects.

# Creating your own Actions

The underlying implementation is a websocket that sends and receives json encoded plaintext messages.

Actions are the requests that you make to the browser in order to automate different tasks.  For instance asking
the browser to navigate to a particular url.  When you construct an Action you need to fill in a step that consists
of the params struct and the reply struct.  In addition you need to specify the API call you are making which is
otherwise known as the MethodName.  These are all defined in the [Devtools Reference](https://chromedevtools.github.io/devtools-protocol/tot).

If your action is going to trigger events that you need to watch for, make sure to include them.

Please refer to example/navigate for a basic example.  This shows an action that consists of a single step and depends on certain
navigation events being fulfilled before the action is considered complete.

# Caveats

Currently there is no code for opening a browser.  There is a start.sh script that shows how to manually start a browser.  
The code will then create a websocket connection for you.

Once a connection is made, you should only run actions against that "frame" in a serial manner.  I haven't tested concurrent access.
It should work, but I cannot guarantee it at the moment.
