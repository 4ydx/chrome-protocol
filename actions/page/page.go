package page

import (
	"github.com/4ydx/cdp/protocol/page"
	"github.com/4ydx/chrome-protocol"
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
func Navigate(frame *cdp.Frame, url string, timeout time.Duration) error {
	return cdp.NewAction(frame,
		GetNavigationEvents(),
		[]cdp.Step{
			cdp.Step{ID: frame.RequestID.GetNext(), Method: page.CommandPageNavigate, Params: &page.NavigateArgs{URL: url}, Reply: &page.NavigateReply{}, Timeout: timeout},
		}).Run()
}
