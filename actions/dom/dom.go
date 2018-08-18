package dom

import (
	"errors"
	"github.com/4ydx/cdp/protocol/dom"
	"github.com/4ydx/cdp/protocol/input"
	"github.com/4ydx/chrome-protocol"
	"log"
	"time"
)

// GetEntireDocument retrieves the root document and all children for the entire page.
func GetEntireDocument(pg *cdp.Page, timeout time.Duration) (*dom.GetFlattenedDocumentReply, error) {
	a0 := cdp.NewAction(pg, []cdp.Event{},
		[]cdp.Step{
			cdp.Step{ID: pg.ID.GetNext(), Method: dom.CommandDOMGetFlattenedDocument, Params: &dom.GetFlattenedDocumentArgs{Depth: -1}, Reply: &dom.GetFlattenedDocumentReply{}, Timeout: timeout},
		})
	err := a0.Run()

	return a0.Steps[0].Reply.(*dom.GetFlattenedDocumentReply), err
}

// FindFirstElementNodeId gets the first element's nodeId using XPath, Css selector, or text matches with the find parameter.
func FindFirstElementNodeId(pg *cdp.Page, find string, timeout time.Duration) (dom.NodeID, error) {
	nodes, err := FindAll(pg, find, timeout)
	if err != nil {
		return 0, err
	}
	if len(nodes) == 0 {
		return 0, errors.New("No element found.")
	}
	target := dom.NodeID(0)
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
func FindAll(pg *cdp.Page, find string, timeout time.Duration) ([]*dom.Node, error) {
	found := make([]*dom.Node, 0)

	doc, err := GetEntireDocument(pg, timeout)
	if err != nil {
		return found, err
	}

	// Make nodeId search request.
	a0 := cdp.NewAction(pg,
		[]cdp.Event{},
		[]cdp.Step{
			cdp.Step{ID: pg.ID.GetNext(), Method: dom.CommandDOMPerformSearch, Params: &dom.PerformSearchArgs{Query: find}, Reply: &dom.PerformSearchReply{}, Timeout: timeout},
		})
	err = a0.Run()
	if err != nil {
		return found, err
	}
	ret := a0.Steps[0].Reply.(*dom.PerformSearchReply)
	if ret.SearchID == "" {
		return found, errors.New("Unexpected empty search id.")
	}
	if ret.ResultCount == 0 {
		return found, errors.New("No nodes found.")
	}

	// Retrieve the NodeIds.
	log.Printf("Using search id %s with result count %d", ret.SearchID, ret.ResultCount)
	params := &dom.GetSearchResultsArgs{
		SearchID:  ret.SearchID,
		FromIndex: 0,
		ToIndex:   ret.ResultCount,
	}
	a1 := cdp.NewAction(pg, []cdp.Event{},
		[]cdp.Step{
			cdp.Step{ID: pg.ID.GetNext(), Method: dom.CommandDOMGetSearchResults, Params: params, Reply: &dom.GetSearchResultsReply{}, Timeout: timeout},
		})
	err = a1.Run()
	if err != nil {
		return found, err
	}

	// Find the matching nodes from the document
	hits := a1.Steps[0].Reply.(*dom.GetSearchResultsReply)
	for _, child := range doc.Nodes {
		for _, id := range hits.NodeIDs {
			if id == child.NodeID {
				found = append(found, &child)
			}
		}
	}

	return found, nil
}

// Focus on the first element node that matches the find parameter.
func Focus(pg *cdp.Page, find string, timeout time.Duration) error {
	target, err := FindFirstElementNodeId(pg, find, timeout)
	if err != nil {
		return err
	}
	a0 := cdp.NewAction(pg, []cdp.Event{},
		[]cdp.Step{
			cdp.Step{ID: pg.ID.GetNext(), Method: dom.CommandDOMFocus, Params: &dom.FocusArgs{NodeID: target}, Reply: &dom.FocusReply{}, Timeout: timeout},
		})
	a0.Run()

	return nil
}

// Click on the first element matching the find parameter.
// nodeID -> DOM.getBoxModel -> Input.dispatchMouseEvent to issue mousedown+mouseup
func Click(pg *cdp.Page, find string, timeout time.Duration) error {
	target, err := FindFirstElementNodeId(pg, find, timeout)
	if err != nil {
		return err
	}
	a0 := cdp.NewAction(pg, []cdp.Event{},
		[]cdp.Step{
			cdp.Step{ID: pg.ID.GetNext(), Method: dom.CommandDOMGetBoxModel, Params: &dom.GetBoxModelArgs{NodeID: target}, Reply: &dom.GetBoxModelReply{}, Timeout: timeout},
		})
	a0.Run()

	// Box is an array of quad vertices, x immediately followed by y for each point, points clock-wise.
	// (0, 1), (2, 3) <- upper edge
	// (4, 5), (6, 7) <- lower edge
	box := a0.Steps[0].Reply.(*dom.GetBoxModelReply).Model.Content
	xMid := (box[2]-box[0])/2 + box[0]
	yMid := (box[5]-box[1])/2 + box[1]

	// Mouse click.
	clicks := []string{"mousePressed", "mouseReleased"}
	for _, click := range clicks {
		a1 := cdp.NewAction(pg, []cdp.Event{},
			[]cdp.Step{
				cdp.Step{ID: pg.ID.GetNext(), Method: input.CommandInputDispatchMouseEvent, Params: &input.DispatchMouseEventArgs{X: xMid, Y: yMid, Button: "left", ClickCount: 1, Type: click}, Reply: &input.DispatchMouseEventReply{}, Timeout: timeout},
			})
		a1.Run()
	}

	return nil
}
