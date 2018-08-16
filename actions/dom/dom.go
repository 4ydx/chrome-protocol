package da

import (
	"github.com/4ydx/cdproto/dom"
	"github.com/4ydx/chrome-protocol"
	"time"
)

// Find finds all nodes using XPath, CSS selector, or text.
func Find(id *cdp.ID, find string, timeout time.Duration, eventCache *cdp.EventCache, actionChan chan<- *cdp.Action, stepChan <-chan bool) *dom.GetSearchResultsReturns {
	// Find nodes on the page if they exist.
	a0 := cdp.NewAction([]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandPerformSearch, Params: &dom.PerformSearchParams{Query: find}, Returns: &dom.PerformSearchReturns{}, Timeout: timeout},
		})
	a0.Run(eventCache, actionChan, stepChan)
	if a0.Steps[0].Returns.(*dom.PerformSearchReturns).SearchID == "" || a0.Steps[0].Returns.(*dom.PerformSearchReturns).ResultCount == 0 {
		return a0.Steps[0].Returns.(*dom.GetSearchResultsReturns)
	}

	// Retrieve the NodeIds.
	a1 := cdp.NewAction([]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandGetSearchResults, Params: &dom.GetSearchResultsParams{}, Returns: &dom.GetSearchResultsReturns{}, Timeout: timeout},
		})
	a1.Run(eventCache, actionChan, stepChan)

	return a1.Steps[0].Returns.(*dom.GetSearchResultsReturns)
}
