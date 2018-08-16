package da

import (
	cd "github.com/4ydx/cdproto/cdp"
	"github.com/4ydx/cdproto/dom"
	"github.com/4ydx/chrome-protocol"
	"log"
	"time"
)

// GetDocument retrieves the root document.
func GetDocument(id *cdp.ID, find string, timeout time.Duration, eventCache *cdp.EventCache, actionChan chan<- *cdp.Action, stepChan <-chan bool) *dom.GetDocumentReturns {
	// Find nodes on the page if they exist.
	a0 := cdp.NewAction([]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandGetDocument, Params: &dom.GetDocumentParams{}, Returns: &dom.GetDocumentReturns{}, Timeout: timeout},
		})
	a0.Run(eventCache, actionChan, stepChan)

	return a0.Steps[0].Returns.(*dom.GetDocumentReturns)
}

// FindAll finds all nodes using XPath, CSS selector, or text.
func FindAll(id *cdp.ID, find string, timeout time.Duration, eventCache *cdp.EventCache, actionChan chan<- *cdp.Action, stepChan <-chan bool) *dom.GetSearchResultsReturns {
	doc := GetDocument(id, find, timeout, eventCache, actionChan, stepChan)
	log.Printf("Doc %+v", doc)

	// Find nodes on the page if they exist.
	a0 := cdp.NewAction([]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandPerformSearch, Params: &dom.PerformSearchParams{Query: find}, Returns: &dom.PerformSearchReturns{}, Timeout: timeout},
		})
	a0.Run(eventCache, actionChan, stepChan)

	ret := a0.Steps[0].Returns.(*dom.PerformSearchReturns)
	if ret.SearchID == "" || ret.ResultCount == 0 {
		return &dom.GetSearchResultsReturns{}
	}

	// Retrieve the NodeIds.
	log.Printf("Using search id %s with result count %d", ret.SearchID, ret.ResultCount)
	params := &dom.GetSearchResultsParams{
		SearchID:  ret.SearchID,
		FromIndex: 0,
		ToIndex:   ret.ResultCount,
	}
	a1 := cdp.NewAction([]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandGetSearchResults, Params: params, Returns: &dom.GetSearchResultsReturns{}, Timeout: timeout},
		})
	a1.Run(eventCache, actionChan, stepChan)

	return a1.Steps[0].Returns.(*dom.GetSearchResultsReturns)
}

// Focus on the node identified by the given nodeId.
func Focus(id *cdp.ID, nodeId cd.NodeID, timeout time.Duration, eventCache *cdp.EventCache, actionChan chan<- *cdp.Action, stepChan <-chan bool) *dom.FocusReturns {
	a0 := cdp.NewAction([]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandFocus, Params: &dom.FocusParams{NodeID: nodeId}, Returns: &dom.FocusReturns{}, Timeout: timeout},
		})
	a0.Run(eventCache, actionChan, stepChan)

	return a0.Steps[0].Returns.(*dom.FocusReturns)
}
