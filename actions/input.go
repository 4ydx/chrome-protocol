package actions

import (
	"github.com/4ydx/cdp/protocol/dom"
	"github.com/4ydx/cdp/protocol/input"
	"github.com/4ydx/chrome-protocol"
	"log"
	"time"
)

// GetWindowsVirtualKeyCode returns the known javascript value for a given modifier key.
// The meta/command modifier is not currently supported as it appear to vary across implementations.
func GetWindowsVirtualKeyCode(modifiers int) int {
	windowsVirtualKeyCode := 0
	switch modifiers {
	case 1:
		// alt
		windowsVirtualKeyCode = 18
	case 2:
		// ctrl
		windowsVirtualKeyCode = 17
	case 4:
		// meta/command - not standard
		windowsVirtualKeyCode = 0
	case 8:
		// shift
		windowsVirtualKeyCode = 16
	}
	return windowsVirtualKeyCode
}

// Fill on the first element node that matches the find parameter.  dom.Focus can be called in order to focus an element in order to fill it.
func Fill(frame *cdp.Frame, find, fill string, timeout time.Duration) error {
	if err := Focus(frame, find, timeout); err != nil {
		return err
	}
	for _, key := range fill {
		err := cdp.NewAction(
			[]cdp.Event{},
			[]cdp.Command{
				cdp.Command{ID: frame.RequestID.GetNext(), Method: input.CommandInputDispatchKeyEvent, Params: &input.DispatchKeyEventArgs{Type: "char", Text: string(key)}, Reply: &input.DispatchKeyEventReply{}, Timeout: timeout},
			}).Run(frame)
		if err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}

// Clear clears out the value attribute of the found element.
func Clear(frame *cdp.Frame, find string, timeout time.Duration) error {
	nodeID, err := FindFirstElementNodeID(frame, find, timeout)
	if err != nil {
		return err
	}
	return SetAttributeValue(frame, nodeID, "value", "", timeout)
}

// KeyDown sends a keydown request to the server.
func KeyDown(frame *cdp.Frame, modifiers int, timeout time.Duration) error {
	windowsVirtualKeyCode := GetWindowsVirtualKeyCode(modifiers)
	err := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: input.CommandInputDispatchKeyEvent, Params: &input.DispatchKeyEventArgs{
				Modifiers:             modifiers,
				Type:                  "keyDown",
				WindowsVirtualKeyCode: windowsVirtualKeyCode,
			}, Reply: &input.DispatchKeyEventReply{}, Timeout: timeout},
		}).Run(frame)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

// MouseScroll scrolls the mouse the given amount.
func MouseScroll(frame *cdp.Frame, deltaX, deltaY float64, timeout time.Duration) error {
	nodes, err := FindAll(frame, "body", timeout)
	if err != nil {
		log.Print(err)
		return err
	}
	nodeID := dom.NodeID(0)
	for _, n := range nodes {
		if n.NodeName == "BODY" {
			nodeID = n.NodeID
		}
	}
	a0 := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: dom.CommandDOMGetBoxModel, Params: &dom.GetBoxModelArgs{NodeID: nodeID}, Reply: &dom.GetBoxModelReply{}, Timeout: timeout},
		})
	err = a0.Run(frame)
	if err != nil {
		log.Print(err)
		return err
	}

	// Box is an array of quad vertices, x immediately followed by y for each point, points clock-wise.
	// (0, 1), (2, 3) <- upper edge
	// (4, 5), (6, 7) <- lower edge
	box := a0.Commands[0].Reply.(*dom.GetBoxModelReply).Model.Content
	xMid := (box[2]-box[0])/2 + box[0]
	yMid := (box[5]-box[1])/2 + box[1]

	// Null values are omited right now in the generated code.  Regardless this command requires both values.
	if deltaX == 0 {
		deltaX = 0.000001
	}
	if deltaY == 0 {
		deltaY = 0.000001
	}
	err = cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: input.CommandInputDispatchMouseEvent, Params: &input.DispatchMouseEventArgs{
				X:          xMid,
				Y:          yMid,
				Button:     "middle",
				ClickCount: 1,
				Type:       "mouseWheel",
				DeltaX:     deltaX,
				DeltaY:     deltaY,
			}, Reply: &input.DispatchMouseEventReply{}, Timeout: timeout},
		}).Run(frame)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
