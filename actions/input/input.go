package input

import (
	"github.com/4ydx/cdp/protocol/input"
	"github.com/4ydx/chrome-protocol"
	"time"
)

// Fill on the first element node that matches the find parameter.
func Fill(pg *cdp.Page, fill string, timeout time.Duration) error {
	for _, key := range fill {
		cdp.NewAction(pg, []cdp.Event{},
			[]cdp.Step{
				cdp.Step{ID: pg.ID.GetNext(), Method: input.CommandInputDispatchKeyEvent, Params: &input.DispatchKeyEventArgs{Type: "char", Text: string(key)}, Reply: &input.DispatchKeyEventReply{}, Timeout: timeout},
			}).Run()
	}
	return nil
}
