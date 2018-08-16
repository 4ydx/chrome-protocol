package da

import (
	//"fmt"
	//"github.com/4ydx/cdproto"
	cd "github.com/4ydx/cdproto/cdp"
	"github.com/4ydx/cdproto/dom"
	"github.com/4ydx/chrome-protocol"
	"log"
	"time"
)

// GetDocument retrieves the root document.
func GetDocument(id *cdp.ID, find string, timeout time.Duration, actionCache *cdp.ActionCache, actionChan chan<- *cdp.Action, stepChan <-chan bool) *dom.GetDocumentReturns {
	// Find nodes on the page if they exist.
	a0 := cdp.NewAction([]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandGetDocument, Params: &dom.GetDocumentParams{Depth: -1}, Returns: &dom.GetDocumentReturns{}, Timeout: timeout},
		})
	a0.Run(actionCache, actionChan, stepChan)

	return a0.Steps[0].Returns.(*dom.GetDocumentReturns)
}

// FindAll finds all nodes using XPath, CSS selector, or text.
func FindAll(id *cdp.ID, find string, timeout time.Duration, actionCache *cdp.ActionCache, actionChan chan<- *cdp.Action, stepChan <-chan bool) *dom.GetSearchResultsReturns {
	doc := GetDocument(id, find, timeout, actionCache, actionChan, stepChan)
	log.Printf("Doc %+v", doc)

	a0 := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandPerformSearch, Params: &dom.PerformSearchParams{Query: find}, Returns: &dom.PerformSearchReturns{}, Timeout: timeout},
		})
	a0.Run(actionCache, actionChan, stepChan)

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
	a1.Run(actionCache, actionChan, stepChan)

	return a1.Steps[0].Returns.(*dom.GetSearchResultsReturns)
}

// Focus on the node identified by the given nodeId.
func Focus(id *cdp.ID, nodeId cd.NodeID, timeout time.Duration, actionCache *cdp.ActionCache, actionChan chan<- *cdp.Action, stepChan <-chan bool) *dom.FocusReturns {
	a0 := cdp.NewAction([]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandFocus, Params: &dom.FocusParams{NodeID: nodeId}, Returns: &dom.FocusReturns{}, Timeout: timeout},
		})
	a0.Run(actionCache, actionChan, stepChan)

	return a0.Steps[0].Returns.(*dom.FocusReturns)
}
