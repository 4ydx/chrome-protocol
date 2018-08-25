package actions

import (
	"github.com/4ydx/cdp/protocol/runtime"
	"github.com/4ydx/chrome-protocol"
	"log"
	"time"
)

// Evaluate runs the javascript expression in the current frame's context.
func Evaluate(frame *cdp.Frame, expression string, timeout time.Duration) (*runtime.EvaluateReply, error) {
	action := cdp.NewAction(frame,
		[]cdp.Event{},
		[]cdp.Step{
			cdp.Step{ID: frame.RequestID.GetNext(), Method: runtime.CommandRuntimeEvaluate, Params: &runtime.EvaluateArgs{Expression: expression, Silent: false}, Reply: &runtime.EvaluateReply{}, Timeout: timeout},
		})
	err := action.Run()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return action.Steps[0].Reply.(*runtime.EvaluateReply), nil
}
