package actions

import (
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
		err := cdp.NewAction(frame,
			[]cdp.Event{},
			[]cdp.Command{
				cdp.Command{ID: frame.RequestID.GetNext(), Method: input.CommandInputDispatchKeyEvent, Params: &input.DispatchKeyEventArgs{Type: "char", Text: string(key)}, Reply: &input.DispatchKeyEventReply{}, Timeout: timeout},
			}).Run()
		if err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}

// KeyDown sends a keydown request to the server.
func KeyDown(frame *cdp.Frame, modifiers int, timeout time.Duration) error {
	windowsVirtualKeyCode := GetWindowsVirtualKeyCode(modifiers)
	err := cdp.NewAction(frame,
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: input.CommandInputDispatchKeyEvent, Params: &input.DispatchKeyEventArgs{
				Modifiers: modifiers,
				Type:      "keyDown",
				WindowsVirtualKeyCode: windowsVirtualKeyCode,
			}, Reply: &input.DispatchKeyEventReply{}, Timeout: timeout},
		}).Run()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
