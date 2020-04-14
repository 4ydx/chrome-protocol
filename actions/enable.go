package actions

import (
	"github.com/4ydx/cdp/protocol/css"
	"github.com/4ydx/cdp/protocol/dom"
	"github.com/4ydx/cdp/protocol/indexeddb"
	"github.com/4ydx/cdp/protocol/inspector"
	"github.com/4ydx/cdp/protocol/log"
	"github.com/4ydx/cdp/protocol/network"
	"github.com/4ydx/cdp/protocol/page"
	"github.com/4ydx/cdp/protocol/runtime"
	"github.com/4ydx/chrome-protocol"
	"time"
)

// EnableAll tells the server to send all event values across the websocket.
func EnableAll(frame *cdp.Frame, timeout time.Duration) error {
	// Order is important.  Dom should come first.
	err := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMEnable, Params: &dom.EnableArgs{}, Reply: &dom.EnableReply{}, Timeout: timeout},
			cdp.Command{ID: frame.RequestID.GetNext(), Method: css.CommandCSSEnable, Params: &css.EnableArgs{}, Reply: &css.EnableReply{}, Timeout: timeout},
			cdp.Command{ID: frame.RequestID.GetNext(), Method: indexeddb.CommandIndexedDBEnable, Params: &indexeddb.EnableArgs{}, Reply: &indexeddb.EnableReply{}, Timeout: timeout},
			cdp.Command{ID: frame.RequestID.GetNext(), Method: inspector.CommandInspectorEnable, Params: &inspector.EnableArgs{}, Reply: &inspector.EnableReply{}, Timeout: timeout},
			cdp.Command{ID: frame.RequestID.GetNext(), Method: log.CommandLogEnable, Params: &log.EnableArgs{}, Reply: &log.EnableReply{}, Timeout: timeout},
			cdp.Command{ID: frame.RequestID.GetNext(), Method: network.CommandNetworkEnable, Params: &network.EnableArgs{}, Reply: &network.EnableReply{}, Timeout: timeout},
			cdp.Command{ID: frame.RequestID.GetNext(), Method: page.CommandPageEnable, Params: &page.EnableArgs{}, Reply: &page.EnableReply{}, Timeout: timeout},
			cdp.Command{ID: frame.RequestID.GetNext(), Method: runtime.CommandRuntimeEnable, Params: &runtime.EnableArgs{}, Reply: &runtime.EnableReply{}, Timeout: timeout},
		}).Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
	}
	return err
}

// EnableDom tells the server to send the dom event values across the websocket.
func EnableDom(frame *cdp.Frame, timeout time.Duration) error {
	err := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMEnable, Params: &dom.EnableArgs{}, Reply: &dom.EnableReply{}, Timeout: timeout},
		}).Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
	}
	return err
}

// EnableLog tells the server to send the log event values across the websocket.
func EnableLog(frame *cdp.Frame, timeout time.Duration) error {
	err := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: log.CommandLogEnable, Params: &log.EnableArgs{}, Reply: &log.EnableReply{}, Timeout: timeout},
		}).Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
	}
	return err
}

// EnablePage tells the server to send the page event values across the websocket.
func EnablePage(frame *cdp.Frame, timeout time.Duration) error {
	err := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: page.CommandPageEnable, Params: &page.EnableArgs{}, Reply: &page.EnableReply{}, Timeout: timeout},
		}).Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
	}
	return err
}

// EnableNetwork tells the server to send the network event values across the websocket.
func EnableNetwork(frame *cdp.Frame, timeout time.Duration) error {
	err := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: network.CommandNetworkEnable, Params: &network.EnableArgs{}, Reply: &network.EnableReply{}, Timeout: timeout},
		}).Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
	}
	return err
}

// EnableRuntime tells the server to send the runtime event values across the websocket.
func EnableRuntime(frame *cdp.Frame, timeout time.Duration) error {
	err := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: runtime.CommandRuntimeEnable, Params: &runtime.EnableArgs{}, Reply: &runtime.EnableReply{}, Timeout: timeout},
		}).Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
	}
	return err
}

// EnableCSS tells the server to send the css event values across the websocket.
func EnableCSS(frame *cdp.Frame, timeout time.Duration) error {
	err := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: css.CommandCSSEnable, Params: &css.EnableArgs{}, Reply: &css.EnableReply{}, Timeout: timeout},
		}).Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
	}
	return err
}

// EnableIndexedDB tells the server to send the indexeddb event values across the websocket.
func EnableIndexedDB(frame *cdp.Frame, timeout time.Duration) error {
	err := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: indexeddb.CommandIndexedDBEnable, Params: &indexeddb.EnableArgs{}, Reply: &indexeddb.EnableReply{}, Timeout: timeout},
		}).Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
	}
	return err
}
