package da

import (
	"errors"
	cd "github.com/4ydx/cdproto/cdp"
	"github.com/4ydx/cdproto/dom"
	"github.com/4ydx/chrome-protocol"
	"log"
	"time"
)

// GetEntireDocument retrieves the root document and all children for the entire page.
func GetEntireDocument(id *cdp.ID, timeout time.Duration) (*dom.GetFlattenedDocumentReturns, error) {
	a0 := cdp.NewAction([]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandGetFlattenedDocument, Params: &dom.GetFlattenedDocumentParams{Depth: -1}, Returns: &dom.GetFlattenedDocumentReturns{}, Timeout: timeout},
		})
	err := a0.Run()

	return a0.Steps[0].Returns.(*dom.GetFlattenedDocumentReturns), err
}

// FindAll finds all nodes using XPath, CSS selector, or text.
func FindAll(id *cdp.ID, find string, timeout time.Duration) (*dom.GetFlattenedDocumentReturns, *dom.GetSearchResultsReturns, error) {
	doc, err := GetEntireDocument(id, timeout)
	if err != nil {
		return doc, &dom.GetSearchResultsReturns{}, err
	}

	a0 := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandPerformSearch, Params: &dom.PerformSearchParams{Query: find}, Returns: &dom.PerformSearchReturns{}, Timeout: timeout},
		})
	err = a0.Run()
	if err != nil {
		return doc, &dom.GetSearchResultsReturns{}, err
	}

	ret := a0.Steps[0].Returns.(*dom.PerformSearchReturns)
	if ret.SearchID == "" || ret.ResultCount == 0 {
		return doc, &dom.GetSearchResultsReturns{}, nil
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
	err = a1.Run()
	if err != nil {
		return doc, &dom.GetSearchResultsReturns{}, err
	}

	return doc, a1.Steps[0].Returns.(*dom.GetSearchResultsReturns), nil
}

// Focus on the first element node that matches the find parameter.
func Focus(id *cdp.ID, find string, timeout time.Duration) error {
	doc, hits, err := FindAll(id, find, timeout)
	if err != nil {
		return err
	}
	target := cd.NodeID(0)
	for _, child := range doc.Nodes {
		for _, id := range hits.NodeIds {
			if id == child.NodeID && child.NodeType == 1 {
				target = id
			}
		}
	}
	if target == 0 {
		return errors.New("No element found.")
	}
	a0 := cdp.NewAction([]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandFocus, Params: &dom.FocusParams{NodeID: target}, Returns: &dom.FocusReturns{}, Timeout: timeout},
		})
	a0.Run()

	return nil
}
