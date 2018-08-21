package input

import (
	"github.com/4ydx/cdp/protocol/input"
	"github.com/4ydx/chrome-protocol"
	"log"
	"time"
)

// Fill on the first element node that matches the find parameter.
func Fill(frame *cdp.Frame, fill string, timeout time.Duration) error {
	for _, key := range fill {
		err := cdp.NewAction(frame, []cdp.Event{},
			[]cdp.Step{
				cdp.Step{ID: frame.RequestID.GetNext(), Method: input.CommandInputDispatchKeyEvent, Params: &input.DispatchKeyEventArgs{Type: "char", Text: string(key)}, Reply: &input.DispatchKeyEventReply{}, Timeout: timeout},
			}).Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
