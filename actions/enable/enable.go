package enable

import (
	"github.com/4ydx/cdp/protocol/css"
	"github.com/4ydx/cdp/protocol/dom"
	"github.com/4ydx/cdp/protocol/inspector"
	"github.com/4ydx/cdp/protocol/log"
	"github.com/4ydx/cdp/protocol/network"
	"github.com/4ydx/cdp/protocol/page"
	"github.com/4ydx/cdp/protocol/runtime"
	"github.com/4ydx/chrome-protocol"
	"time"
)

func All(frame *cdp.Frame, timeout time.Duration) error {
	// Order is important.  Dom should come first.
	return cdp.NewAction(frame, []cdp.Event{}, []cdp.Step{
		cdp.Step{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMEnable, Params: &dom.EnableArgs{}, Reply: &dom.EnableReply{}, Timeout: timeout},
		cdp.Step{ID: frame.RequestID.GetNext(), Method: css.CommandCSSEnable, Params: &css.EnableArgs{}, Reply: &css.EnableReply{}, Timeout: timeout},
		cdp.Step{ID: frame.RequestID.GetNext(), Method: inspector.CommandInspectorEnable, Params: &inspector.EnableArgs{}, Reply: &inspector.EnableReply{}, Timeout: timeout},
		cdp.Step{ID: frame.RequestID.GetNext(), Method: log.CommandLogEnable, Params: &log.EnableArgs{}, Reply: &log.EnableReply{}, Timeout: timeout},
		cdp.Step{ID: frame.RequestID.GetNext(), Method: network.CommandNetworkEnable, Params: &network.EnableArgs{}, Reply: &network.EnableReply{}, Timeout: timeout},
		cdp.Step{ID: frame.RequestID.GetNext(), Method: page.CommandPageEnable, Params: &page.EnableArgs{}, Reply: &page.EnableReply{}, Timeout: timeout},
		cdp.Step{ID: frame.RequestID.GetNext(), Method: runtime.CommandRuntimeEnable, Params: &runtime.EnableArgs{}, Reply: &runtime.EnableReply{}, Timeout: timeout},
	}).Run()
}

func Dom(frame *cdp.Frame, timeout time.Duration) error {
	return cdp.NewAction(frame, []cdp.Event{}, []cdp.Step{
		cdp.Step{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMEnable, Params: &dom.EnableArgs{}, Reply: &dom.EnableReply{}, Timeout: timeout},
	}).Run()
}

func Page(frame *cdp.Frame, timeout time.Duration) error {
	return cdp.NewAction(frame, []cdp.Event{}, []cdp.Step{
		cdp.Step{ID: frame.RequestID.GetNext(), Method: page.CommandPageEnable, Params: &page.EnableArgs{}, Reply: &page.EnableReply{}, Timeout: timeout},
	}).Run()
}

func Network(frame *cdp.Frame, timeout time.Duration) error {
	return cdp.NewAction(frame, []cdp.Event{}, []cdp.Step{
		cdp.Step{ID: frame.RequestID.GetNext(), Method: network.CommandNetworkEnable, Params: &network.EnableArgs{}, Reply: &network.EnableReply{}, Timeout: timeout},
	}).Run()
}
