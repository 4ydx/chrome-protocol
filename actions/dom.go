package actions

import (
	"errors"
	"github.com/4ydx/cdp/protocol/dom"
	"github.com/4ydx/cdp/protocol/input"
	"github.com/4ydx/chrome-protocol"
	"time"
)

// GetEntireDocument retrieves the root document and all children for the entire page.
func GetEntireDocument(frame *cdp.Frame, timeout time.Duration) (*dom.GetFlattenedDocumentReply, error) {
	frameDOM := frame.GetDOM()
	if frameDOM != nil && len(frameDOM.Nodes) > 0 {
		frame.Browser.Log.Print("Using cached Frame DOM.")
		return frameDOM, nil
	}
	a0 := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMGetFlattenedDocument, Params: &dom.GetFlattenedDocumentArgs{Depth: -1}, Reply: &dom.GetFlattenedDocumentReply{}, Timeout: timeout},
		})
	err := a0.Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
		return nil, err
	}
	frame.SetDOM(a0.Commands[0].Reply.(*dom.GetFlattenedDocumentReply))

	return a0.Commands[0].Reply.(*dom.GetFlattenedDocumentReply), err
}

// FindAll finds all nodes using XPath, CSS selector, or text.
func FindAll(frame *cdp.Frame, find string, timeout time.Duration) ([]dom.Node, error) {
	found := make([]dom.Node, 0)

	doc, err := GetEntireDocument(frame, timeout)
	if err != nil {
		frame.Browser.Log.Print(err)
		return found, err
	}

	// Make nodeId search request.
	a0 := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMPerformSearch, Params: &dom.PerformSearchArgs{Query: find}, Reply: &dom.PerformSearchReply{}, Timeout: timeout},
		})
	err = a0.Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
		return found, err
	}
	ret := a0.Commands[0].Reply.(*dom.PerformSearchReply)
	if ret.SearchID == "" {
		err := errors.New("unexpected empty search id")
		frame.Browser.Log.Print(err)
		return found, err
	}
	if ret.ResultCount == 0 {
		err := errors.New("no nodes found")
		frame.Browser.Log.Print(err)
		return found, err
	}

	// Retrieve the NodeIds.
	params := &dom.GetSearchResultsArgs{
		SearchID:  ret.SearchID,
		FromIndex: 0,
		ToIndex:   ret.ResultCount,
	}
	a1 := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMGetSearchResults, Params: params, Reply: &dom.GetSearchResultsReply{}, Timeout: timeout},
		})
	err = a1.Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
		return found, err
	}

	// Find the matching nodes from the document
	hits := a1.Commands[0].Reply.(*dom.GetSearchResultsReply)
	for _, child := range doc.Nodes {
		for _, id := range hits.NodeIDs {
			if id == child.NodeID {
				found = append(found, child)
			}
		}
	}
	return found, nil
}

// FindFirstElementNodeID gets the first element's nodeId using XPath, Css selector, or text matches with the find parameter.
func FindFirstElementNodeID(frame *cdp.Frame, find string, timeout time.Duration) (dom.NodeID, error) {
	nodes, err := FindAll(frame, find, timeout)
	if err != nil {
		frame.Browser.Log.Print(err)
		return 0, err
	}
	if len(nodes) == 0 {
		err := errors.New("no element found")
		frame.Browser.Log.Print(err)
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
		frame.Browser.Log.Print(err)
		return 0, err
	}
	return target, nil
}

// Focus on the first element node that matches the find parameter.
func Focus(frame *cdp.Frame, find string, timeout time.Duration) error {
	target, err := FindFirstElementNodeID(frame, find, timeout)
	if err != nil {
		frame.Browser.Log.Print(err)
		return err
	}
	err = cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMFocus, Params: &dom.FocusArgs{NodeID: target}, Reply: &dom.FocusReply{}, Timeout: timeout},
		}).Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
		return err
	}
	return nil
}

// Click on the first element matching the find parameter.
// Any events that need to be tracked as a result of the click must be included.
// This will insure that the click action waits until required events are fired.
func Click(frame *cdp.Frame, find string, events []cdp.Event, timeout time.Duration) ([]cdp.Event, error) {
	return ClickWithModifiers(frame, find, 0, events, timeout)
}

// ClickWithModifiers clicks on a found element using the specified key modifier values.
func ClickWithModifiers(frame *cdp.Frame, find string, modifiers int, events []cdp.Event, timeout time.Duration) ([]cdp.Event, error) {
	target, err := FindFirstElementNodeID(frame, find, timeout)
	if err != nil {
		frame.Browser.Log.Print(err)
		return events, err
	}
	return ClickNodeID(frame, target, modifiers, events, timeout)
}

// ClickNodeID clicks on the element identified by the given dom.NodeID value.
func ClickNodeID(frame *cdp.Frame, nodeID dom.NodeID, modifiers int, events []cdp.Event, timeout time.Duration) ([]cdp.Event, error) {
	a0 := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMGetBoxModel, Params: &dom.GetBoxModelArgs{NodeID: nodeID}, Reply: &dom.GetBoxModelReply{}, Timeout: timeout},
		})
	err := a0.Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
		return events, err
	}

	// Box is an array of quad vertices, x immediately followed by y for each point, points clock-wise.
	// (0, 1), (2, 3) <- upper edge
	// (4, 5), (6, 7) <- lower edge
	box := a0.Commands[0].Reply.(*dom.GetBoxModelReply).Model.Content
	xMid := (box[2]-box[0])/2 + box[0]
	yMid := (box[5]-box[1])/2 + box[1]

	// Mouse click.
	left := input.MouseButtonLeft
	err = cdp.NewAction(
		events,
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: input.CommandInputDispatchMouseEvent, Params: &input.DispatchMouseEventArgs{
				Modifiers:  modifiers,
				X:          xMid,
				Y:          yMid,
				Button:     &left,
				ClickCount: 1,
				Type:       "mousePressed",
			}, Reply: &input.DispatchMouseEventReply{}, Timeout: timeout},
			cdp.Command{ID: frame.RequestID.GetNext(), Method: input.CommandInputDispatchMouseEvent, Params: &input.DispatchMouseEventArgs{
				Modifiers:  modifiers,
				X:          xMid,
				Y:          yMid,
				Button:     &left,
				ClickCount: 1,
				Type:       "mouseReleased",
			}, Reply: &input.DispatchMouseEventReply{}, Timeout: timeout},
		}).Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
		return events, err
	}
	return events, nil
}

// Children of the first element node that matches the find parameter.  If the frame.DOM object already has the data, this call will do nothing.  Otherwise, it should trigger DOM.setChildNodes events.
// NOTE: It appears that before this action will be completed (before the reply is received), if the server has not yet sent any/some of the child nodes of the given nodeID, then it will send those to the client
//       as DOM.setChildNodes events.  We do not need to pick those up here since there is a method in the websocket loop of github.com/4ydx/chrome-protocol that watches for such events and updates the DOM object.
//       In fact, I don't know how many of those events might be fired and an action's event slice isn't designed to handle multiples of the same event type.
func Children(frame *cdp.Frame, nodeID dom.NodeID, timeout time.Duration) error {
	err := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMRequestChildNodes, Params: &dom.RequestChildNodesArgs{NodeID: nodeID, Depth: -1}, Reply: &dom.RequestChildNodesReply{}, Timeout: timeout},
		}).Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
		return err
	}
	return nil
}

// SetAttributeValue sets the value of the given attribute of the given nodeID to the given value.
func SetAttributeValue(frame *cdp.Frame, nodeID dom.NodeID, name, value string, timeout time.Duration) error {
	err := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMSetAttributeValue, Params: &dom.SetAttributeValueArgs{NodeID: nodeID, Name: name, Value: value}, Reply: &dom.SetAttributeValueReply{}, Timeout: timeout},
		}).Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
		return err
	}
	return nil
}
