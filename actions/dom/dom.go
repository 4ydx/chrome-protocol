package dom

import (
	"errors"
	cd "github.com/4ydx/cdproto/cdp"
	"github.com/4ydx/cdproto/dom"
	"github.com/4ydx/cdproto/input"
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

// FindFirstElementNodeId gets the first element's nodeId using XPath, Css selector, or text matches with the find parameter.
func FindFirstElementNodeId(id *cdp.ID, find string, timeout time.Duration) (cd.NodeID, error) {
	nodes, err := FindAll(id, find, timeout)
	if err != nil {
		return 0, err
	}
	if len(nodes) == 0 {
		return 0, errors.New("No element found.")
	}
	target := cd.NodeID(0)
	for _, child := range nodes {
		if child.NodeType == 1 {
			// Is element node.
			target = child.NodeID
			break
		}
	}
	return target, nil
}

// FindAll finds all nodes using XPath, CSS selector, or text.
func FindAll(id *cdp.ID, find string, timeout time.Duration) ([]*cd.Node, error) {
	found := make([]*cd.Node, 0)

	doc, err := GetEntireDocument(id, timeout)
	if err != nil {
		return found, err
	}

	// Make nodeId search request.
	a0 := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandPerformSearch, Params: &dom.PerformSearchParams{Query: find}, Returns: &dom.PerformSearchReturns{}, Timeout: timeout},
		})
	err = a0.Run()
	if err != nil {
		return found, err
	}
	ret := a0.Steps[0].Returns.(*dom.PerformSearchReturns)
	if ret.SearchID == "" {
		return found, errors.New("Unexpected empty search id.")
	}
	if ret.ResultCount == 0 {
		return found, errors.New("No nodes found.")
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
		return found, err
	}

	// Find the matching nodes from the document
	hits := a1.Steps[0].Returns.(*dom.GetSearchResultsReturns)
	for _, child := range doc.Nodes {
		for _, id := range hits.NodeIds {
			if id == child.NodeID {
				found = append(found, child)
			}
		}
	}

	return found, nil
}

// Focus on the first element node that matches the find parameter.
func Focus(id *cdp.ID, find string, timeout time.Duration) error {
	target, err := FindFirstElementNodeId(id, find, timeout)
	if err != nil {
		return err
	}
	a0 := cdp.NewAction([]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandFocus, Params: &dom.FocusParams{NodeID: target}, Returns: &dom.FocusReturns{}, Timeout: timeout},
		})
	a0.Run()

	return nil
}

// Click on the first element matching the find parameter.
// nodeID -> DOM.getBoxModel -> Input.dispatchMouseEvent to issue mousedown+mouseup
func Click(id *cdp.ID, find string, timeout time.Duration) error {
	target, err := FindFirstElementNodeId(id, find, timeout)
	if err != nil {
		return err
	}
	a0 := cdp.NewAction([]cdp.Event{},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: dom.CommandGetBoxModel, Params: &dom.GetBoxModelParams{NodeID: target}, Returns: &dom.GetBoxModelReturns{}, Timeout: timeout},
		})
	a0.Run()

	// Box is an array of quad vertices, x immediately followed by y for each point, points clock-wise.
	// (0, 1), (2, 3) <- upper edge
	// (4, 5), (6, 7) <- lower edge
	box := a0.Steps[0].Returns.(*dom.GetBoxModelReturns).Model.Content
	xMid := (box[2]-box[0])/2 + box[0]
	yMid := (box[5]-box[1])/2 + box[1]

	// Mouse click.
	clicks := []input.MouseType{input.MousePressed, input.MouseReleased}
	for _, click := range clicks {
		a1 := cdp.NewAction([]cdp.Event{},
			[]cdp.Step{
				cdp.Step{Id: id.GetNext(), Method: input.CommandDispatchMouseEvent, Params: &input.DispatchMouseEventParams{X: xMid, Y: yMid, Button: input.ButtonLeft, ClickCount: 1, Type: click}, Returns: &input.DispatchMouseEventReturns{}, Timeout: timeout},
			})
		a1.Run()
	}

	return nil
}
