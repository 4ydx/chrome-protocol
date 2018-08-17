package input

import (
	"github.com/4ydx/cdproto/input"
	"github.com/4ydx/chrome-protocol"
	"time"
)

// Fill on the first element node that matches the find parameter.
func Fill(id *cdp.ID, fill string, timeout time.Duration) error {
	for _, key := range fill {
		a0 := cdp.NewAction([]cdp.Event{},
			[]cdp.Step{
				cdp.Step{Id: id.GetNext(), Method: input.CommandDispatchKeyEvent, Params: &input.DispatchKeyEventParams{Type: "char", Text: string(key)}, Returns: &input.DispatchKeyEventReturns{}, Timeout: timeout},
			})
		a0.Run()
		a0.Log()
	}
	return nil
}
