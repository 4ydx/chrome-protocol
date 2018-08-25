package actions

import (
	"errors"
	"github.com/4ydx/cdp/protocol/dom"
	"github.com/4ydx/cdp/protocol/input"
	"github.com/4ydx/chrome-protocol"
	"log"
	"time"
)

// GetEntireDocument retrieves the root document and all children for the entire page.
func GetEntireDocument(frame *cdp.Frame, timeout time.Duration) (*dom.GetFlattenedDocumentReply, error) {
	a0 := cdp.NewAction(frame, []cdp.Event{},
		[]cdp.Step{
			cdp.Step{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMGetFlattenedDocument, Params: &dom.GetFlattenedDocumentArgs{Depth: -1}, Reply: &dom.GetFlattenedDocumentReply{}, Timeout: timeout},
		})
	err := a0.Run()
	if err != nil {
		log.Print(err)
	}
	return a0.Steps[0].Reply.(*dom.GetFlattenedDocumentReply), err
}

// FindFirstElementNodeID gets the first element's nodeId using XPath, Css selector, or text matches with the find parameter.
func FindFirstElementNodeID(frame *cdp.Frame, find string, timeout time.Duration) (dom.NodeID, error) {
	nodes, err := FindAll(frame, find, timeout)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	if len(nodes) == 0 {
		err := errors.New("no element found")
		log.Print(err)
		return 0, err
	}
	target := dom.NodeID(0)
	for _, child := range nodes {
		if child.NodeType == 1 {
			// Is node of type 'element'.
			target = child.NodeID
			break
		}
	}
	if target == 0 {
		err := errors.New("no element (NodeType 1) found within matching nodes")
		log.Print(err)
		return 0, err
	}
	return target, nil
}

// FindAll finds all nodes using XPath, CSS selector, or text.
func FindAll(frame *cdp.Frame, find string, timeout time.Duration) ([]dom.Node, error) {
	found := make([]dom.Node, 0)

	doc, err := GetEntireDocument(frame, timeout)
	if err != nil {
		log.Print(err)
		return found, err
	}

	// Make nodeId search request.
	a0 := cdp.NewAction(frame,
		[]cdp.Event{},
		[]cdp.Step{
			cdp.Step{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMPerformSearch, Params: &dom.PerformSearchArgs{Query: find}, Reply: &dom.PerformSearchReply{}, Timeout: timeout},
		})
	err = a0.Run()
	if err != nil {
		log.Print(err)
		return found, err
	}
	ret := a0.Steps[0].Reply.(*dom.PerformSearchReply)
	if ret.SearchID == "" {
		err := errors.New("unexpected empty search id")
		log.Print(err)
		return found, err
	}
	if ret.ResultCount == 0 {
		err := errors.New("no nodes found")
		log.Print(err)
		return found, err
	}

	// Retrieve the NodeIds.
	params := &dom.GetSearchResultsArgs{
		SearchID:  ret.SearchID,
		FromIndex: 0,
		ToIndex:   ret.ResultCount,
	}
	a1 := cdp.NewAction(frame, []cdp.Event{},
		[]cdp.Step{
			cdp.Step{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMGetSearchResults, Params: params, Reply: &dom.GetSearchResultsReply{}, Timeout: timeout},
		})
	err = a1.Run()
	if err != nil {
		log.Print(err)
		return found, err
	}

	// Find the matching nodes from the document
	hits := a1.Steps[0].Reply.(*dom.GetSearchResultsReply)
	for _, child := range doc.Nodes {
		for _, id := range hits.NodeIDs {
			if id == child.NodeID {
				found = append(found, child)
			}
		}
	}
	return found, nil
}

// Focus on the first element node that matches the find parameter.
func Focus(frame *cdp.Frame, find string, timeout time.Duration) error {
	target, err := FindFirstElementNodeID(frame, find, timeout)
	if err != nil {
		log.Print(err)
		return err
	}
	err = cdp.NewAction(frame, []cdp.Event{},
		[]cdp.Step{
			cdp.Step{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMFocus, Params: &dom.FocusArgs{NodeID: target}, Reply: &dom.FocusReply{}, Timeout: timeout},
		}).Run()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

// Click on the first element matching the find parameter.
// Any events that need to be tracked as a result of the click must be included.
// This will insure that the click action waits until required events are fired.
func Click(frame *cdp.Frame, find string, events []cdp.Event, timeout time.Duration) ([]cdp.Event, error) {
	target, err := FindFirstElementNodeID(frame, find, timeout)
	if err != nil {
		log.Print(err)
		return events, err
	}
	a0 := cdp.NewAction(frame, []cdp.Event{},
		[]cdp.Step{
			cdp.Step{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMGetBoxModel, Params: &dom.GetBoxModelArgs{NodeID: target}, Reply: &dom.GetBoxModelReply{}, Timeout: timeout},
		})
	err = a0.Run()
	if err != nil {
		log.Print(err)
		return events, err
	}

	// Box is an array of quad vertices, x immediately followed by y for each point, points clock-wise.
	// (0, 1), (2, 3) <- upper edge
	// (4, 5), (6, 7) <- lower edge
	box := a0.Steps[0].Reply.(*dom.GetBoxModelReply).Model.Content
	xMid := (box[2]-box[0])/2 + box[0]
	yMid := (box[5]-box[1])/2 + box[1]

	// Mouse click.
	err = cdp.NewAction(frame, events,
		[]cdp.Step{
			cdp.Step{ID: frame.RequestID.GetNext(), Method: input.CommandInputDispatchMouseEvent, Params: &input.DispatchMouseEventArgs{X: xMid, Y: yMid, Button: "left", ClickCount: 1, Type: "mousePressed"}, Reply: &input.DispatchMouseEventReply{}, Timeout: timeout},
			cdp.Step{ID: frame.RequestID.GetNext(), Method: input.CommandInputDispatchMouseEvent, Params: &input.DispatchMouseEventArgs{X: xMid, Y: yMid, Button: "left", ClickCount: 1, Type: "mouseReleased"}, Reply: &input.DispatchMouseEventReply{}, Timeout: timeout},
		}).Run()
	if err != nil {
		log.Print(err)
		return events, err
	}
	return events, nil
}
