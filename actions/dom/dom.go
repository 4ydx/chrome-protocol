package da

import (
	"github.com/4ydx/cdproto/dom"
	"github.com/4ydx/chrome-protocol"
	"time"
)

// Find finds all nodes using XPath, CSS selector, or text.
func Find(id *cdp.ID, find string, timeout time.Duration) *cdp.Action {
	return cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandPerformSearch, Params: &dom.PerformSearchParams{Query: find}, Returns: &dom.PerformSearchReturns{}, Timeout: timeout},
		})
}
