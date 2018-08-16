package da

import (
	cd "github.com/4ydx/cdproto/cdp"
	"github.com/4ydx/cdproto/dom"
	"github.com/4ydx/chrome-protocol"
	"log"
	"time"
)

// GetEntireDocument retrieves the root document and all children for the entire page.
func GetEntireDocument(id *cdp.ID, timeout time.Duration, actionCache *cdp.ActionCache, actionChan chan<- *cdp.Action, stepChan <-chan bool) *dom.GetFlattenedDocumentReturns {
	a0 := cdp.NewAction([]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandGetFlattenedDocument, Params: &dom.GetFlattenedDocumentParams{Depth: -1}, Returns: &dom.GetFlattenedDocumentReturns{}, Timeout: timeout},
		})
	a0.Run(actionCache, actionChan, stepChan)

	return a0.Steps[0].Returns.(*dom.GetFlattenedDocumentReturns)
}

// FindAll finds all nodes using XPath, CSS selector, or text.
func FindAll(id *cdp.ID, find string, timeout time.Duration, actionCache *cdp.ActionCache, actionChan chan<- *cdp.Action, stepChan <-chan bool) (*dom.GetFlattenedDocumentReturns, *dom.GetSearchResultsReturns) {
	doc := GetEntireDocument(id, timeout, actionCache, actionChan, stepChan)

	a0 := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandPerformSearch, Params: &dom.PerformSearchParams{Query: find}, Returns: &dom.PerformSearchReturns{}, Timeout: timeout},
		})
	a0.Run(actionCache, actionChan, stepChan)

	ret := a0.Steps[0].Returns.(*dom.PerformSearchReturns)
	if ret.SearchID == "" || ret.ResultCount == 0 {
		return doc, &dom.GetSearchResultsReturns{}
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

	return doc, a1.Steps[0].Returns.(*dom.GetSearchResultsReturns)
}

// Focus on the node identified by the given nodeId.
func Focus(id *cdp.ID, find string, timeout time.Duration, actionCache *cdp.ActionCache, actionChan chan<- *cdp.Action, stepChan <-chan bool) *dom.FocusReturns {
	doc, hits := FindAll(id, find, timeout, actionCache, actionChan, stepChan)

	target := cd.NodeID(0)
	for _, child := range doc.Nodes {
		for _, id := range hits.NodeIds {
			if id == child.NodeID && child.NodeType == 1 {
				target = id
			}
		}
	}
	a0 := cdp.NewAction([]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandFocus, Params: &dom.FocusParams{NodeID: target}, Returns: &dom.FocusReturns{}, Timeout: timeout},
		})
	a0.Run(actionCache, actionChan, stepChan)

	return a0.Steps[0].Returns.(*dom.FocusReturns)
}