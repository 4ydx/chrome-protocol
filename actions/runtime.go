package actions

import (
	"github.com/4ydx/cdp/protocol"
	"github.com/4ydx/cdp/protocol/runtime"
	"github.com/4ydx/chrome-protocol"
	"log"
	"time"
)

// Evaluate runs the javascript expression in the current frame's context.
func Evaluate(frame *cdp.Frame, expression string, timeout time.Duration) (*runtime.EvaluateReply, error) {
	action := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: runtime.CommandRuntimeEvaluate, Params: &runtime.EvaluateArgs{Expression: expression, Silent: false}, Reply: &runtime.EvaluateReply{}, Timeout: timeout},
		})
	err := action.Run(frame)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return action.Commands[0].Reply.(*runtime.EvaluateReply), nil
}

// GetProperties runs the properties of a given object.
func GetProperties(frame *cdp.Frame, objectID shared.RemoteObjectID, ownProperties, accessorPropertiesOnly bool, timeout time.Duration) (*runtime.GetPropertiesReply, error) {
	args := &runtime.GetPropertiesArgs{
		ObjectID:               objectID,
		OwnProperties:          ownProperties,
		AccessorPropertiesOnly: accessorPropertiesOnly,
	}
	action := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: runtime.CommandRuntimeGetProperties, Params: args, Reply: &runtime.GetPropertiesReply{}, Timeout: timeout},
		})
	err := action.Run(frame)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return action.Commands[0].Reply.(*runtime.GetPropertiesReply), nil
}
