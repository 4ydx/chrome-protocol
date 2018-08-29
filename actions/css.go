package actions

import (
	"fmt"
	"github.com/4ydx/cdp/protocol/css"
	"github.com/4ydx/cdp/protocol/dom"
	"github.com/4ydx/chrome-protocol"
	"log"
	"time"
)

// GetComputedStyleForNode get the computed style for a node.
func GetComputedStyleForNode(frame *cdp.Frame, nodeID dom.NodeID, timeout time.Duration) (*css.GetComputedStyleForNodeReply, error) {
	action := cdp.NewAction(frame,
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: css.CommandCSSGetComputedStyleForNode, Params: &css.GetComputedStyleForNodeArgs{NodeID: nodeID}, Reply: &css.GetComputedStyleForNodeReply{}, Timeout: timeout},
		})
	err := action.Run()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return action.Commands[0].Reply.(*css.GetComputedStyleForNodeReply), nil
}

// WaitForComputedStyle finds the first element on a page by id, css, or xpath and waits until the given css propery is set to the given css value.
// For instance, wait until the cssPropery "display" is set to cssValue "none".  That is, until the found html element disappears from view.
func WaitForComputedStyle(frame *cdp.Frame, find, cssPropery, cssValue string, timeout time.Duration) error {
	nodeID, err := FindFirstElementNodeID(frame, find, timeout)
	if err != nil {
		log.Print(err)
		return err
	}
	until := time.Now().Add(timeout)
	for {
		if time.Now().After(until) {
			log.Print("timeout")
			return fmt.Errorf("timeout")
		}
		style, err := GetComputedStyleForNode(frame, nodeID, timeout)
		if err != nil {
			log.Print(err)
			return err
		}
		visible := true
		for _, s := range style.ComputedStyle {
			if s.Name == cssPropery && s.Value == cssValue {
				visible = false
			}
		}
		if !visible {
			break
		}
	}
	return nil
}
