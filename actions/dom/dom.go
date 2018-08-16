package da

import (
	"github.com/4ydx/cdproto/dom"
	"github.com/4ydx/chrome-protocol"
	"time"
)

// Find finds all nodes using XPath, CSS selector, or text.
// TODO: so how to chain this together???
//       a "step" needs to be able to reliably get a value from a previous step and then apply that value to its own params.
func Find(actions *cdp.Action, actions []*cdp.Action, id *cdp.ID, find string, timeout time.Duration) *cdp.Action {
	act := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id0, Method: dom.CommandPerformSearch, Params: &dom.PerformSearchParams{Query: find}, Returns: searchReturns, Timeout: timeout},
			cdp.Step{Id: id.GetNext(), Method: dom.CommandGetSearchResults, Params: getSearchResults, Returns: &dom.GetSearchResultsReturns{}, Timeout: timeout, PreviousReturns: func() {
				getSearchResults.SearchID = searchResults.SearchID
			}},
		})
}
