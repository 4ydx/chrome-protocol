package ea

import (
	"github.com/4ydx/cdproto/css"
	"github.com/4ydx/cdproto/dom"
	"github.com/4ydx/cdproto/inspector"
	"github.com/4ydx/cdproto/log"
	"github.com/4ydx/cdproto/network"
	"github.com/4ydx/cdproto/page"
	"github.com/4ydx/cdproto/runtime"
	"github.com/4ydx/chrome-protocol"
	"time"
)

func EnableAll(id *cdp.ID, timeout time.Duration) error {
	// Order is important.  Dom should come first.
	return cdp.NewAction([]cdp.Event{}, []cdp.Step{
		cdp.Step{Id: id.GetNext(), Method: dom.CommandEnable, Params: &dom.EnableParams{}, Returns: &dom.EnableReturns{}, Timeout: timeout},

		cdp.Step{Id: id.GetNext(), Method: css.CommandEnable, Params: &css.EnableParams{}, Returns: &css.EnableReturns{}, Timeout: timeout},
		cdp.Step{Id: id.GetNext(), Method: inspector.CommandEnable, Params: &inspector.EnableParams{}, Returns: &inspector.EnableReturns{}, Timeout: timeout},
		cdp.Step{Id: id.GetNext(), Method: log.CommandEnable, Params: &log.EnableParams{}, Returns: &log.EnableReturns{}, Timeout: timeout},
		cdp.Step{Id: id.GetNext(), Method: network.CommandEnable, Params: &network.EnableParams{}, Returns: &network.EnableReturns{}, Timeout: timeout},
		cdp.Step{Id: id.GetNext(), Method: page.CommandEnable, Params: &page.EnableParams{}, Returns: &page.EnableReturns{}, Timeout: timeout},
		cdp.Step{Id: id.GetNext(), Method: runtime.CommandEnable, Params: &runtime.EnableParams{}, Returns: &runtime.EnableReturns{}, Timeout: timeout},
	}).Run()
}

func EnableDom(id *cdp.ID, timeout time.Duration) error {
	return cdp.NewAction([]cdp.Event{}, []cdp.Step{
		cdp.Step{Id: id.GetNext(), Method: dom.CommandEnable, Params: &dom.EnableParams{}, Returns: &dom.EnableReturns{}, Timeout: timeout},
	}).Run()
}

func EnablePage(id *cdp.ID, timeout time.Duration) error {
	return cdp.NewAction([]cdp.Event{}, []cdp.Step{
		cdp.Step{Id: id.GetNext(), Method: page.CommandEnable, Params: &page.EnableParams{}, Returns: &page.EnableReturns{}, Timeout: timeout},
	}).Run()
}
