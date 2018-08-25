package actions

import (
	"bytes"
	"github.com/4ydx/cdp/protocol/page"
	"github.com/4ydx/chrome-protocol"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
	"time"
)

// GetNavigationEvents returns all events that are expected to occur after a page navigation api request is made.
func GetNavigationEvents() []cdp.Event {
	return []cdp.Event{
		cdp.Event{Name: page.EventPageFrameNavigated, Value: &page.FrameNavigatedReply{}, IsRequired: true},
		cdp.Event{Name: page.EventPageFrameStartedLoading, Value: &page.FrameStartedLoadingReply{}, IsRequired: true},
		cdp.Event{Name: page.EventPageFrameStoppedLoading, Value: &page.FrameStoppedLoadingReply{}, IsRequired: true},
	}
}

// Navigate sends the browser to the given URL
func Navigate(frame *cdp.Frame, url string, timeout time.Duration) ([]cdp.Event, error) {
	events := GetNavigationEvents()
	action := cdp.NewAction(frame,
		events,
		[]cdp.Step{
			cdp.Step{ID: frame.RequestID.GetNext(), Method: page.CommandPageNavigate, Params: &page.NavigateArgs{URL: url}, Reply: &page.NavigateReply{}, Timeout: timeout},
		})
	if err := action.Run(); err != nil {
		log.Print(err)
	}
	return events, nil
}

// Screenshot captures a screenshot and saves it to the given destination.
func Screenshot(frame *cdp.Frame, destination, format string, quality int, clip *page.Viewport, timeout time.Duration) error {
	var action *cdp.Action
	if clip != nil {
		action = cdp.NewAction(frame,
			[]cdp.Event{},
			[]cdp.Step{
				cdp.Step{ID: frame.RequestID.GetNext(), Method: page.CommandPageCaptureScreenshot, Params: &page.CaptureScreenshotArgs{Format: format, Clip: clip, Quality: quality}, Reply: &page.CaptureScreenshotReply{}, Timeout: timeout},
			})
	} else {
		action = cdp.NewAction(frame,
			[]cdp.Event{},
			[]cdp.Step{
				cdp.Step{ID: frame.RequestID.GetNext(), Method: page.CommandPageCaptureScreenshot, Params: &page.CaptureScreenshotArgs{Format: format, Quality: quality}, Reply: &page.CaptureScreenshotReply{}, Timeout: timeout},
			})
	}
	if err := action.Run(); err != nil {
		log.Print(err)
		return err
	}

	// Convert to an image.
	src := action.Steps[0].Reply.(*page.CaptureScreenshotReply).Data
	m, _, err := image.Decode(bytes.NewReader(src))
	if err != nil {
		log.Print(err)
		return err
	}

	// Save to destination.
	if !strings.HasSuffix(destination, "."+format) {
		destination = destination + "." + format
	}
	f, err := os.Create(destination)
	if err != nil {
		log.Print(err)
		return err
	}
	defer f.Close()

	if format == "png" {
		png.Encode(f, m)
	} else {
		jpeg.Encode(f, m, nil)
	}
	return nil
}
