package page

import (
	"github.com/4ydx/cdp/protocol/page"
	"github.com/4ydx/chrome-protocol"
	"time"
)

func GetNavigationEvents() []cdp.Event {
	return []cdp.Event{
		cdp.Event{Name: page.EventPageFrameNavigated, Value: &page.FrameNavigatedReply{}, IsRequired: true},
		cdp.Event{Name: page.EventPageFrameStartedLoading, Value: &page.FrameStartedLoadingReply{}, IsRequired: true},
		cdp.Event{Name: page.EventPageFrameStoppedLoading, Value: &page.FrameStoppedLoadingReply{}, IsRequired: true},
	}
}

// Navigate sends the browser to the given URL
func Navigate(pg *cdp.Frame, url string, timeout time.Duration) error {
	return cdp.NewAction(pg,
		GetNavigationEvents(),
		[]cdp.Step{
			cdp.Step{ID: pg.RequestID.GetNext(), Method: page.CommandPageNavigate, Params: &page.NavigateArgs{URL: url}, Reply: &page.NavigateReply{}, Timeout: timeout},
		}).Run()
}
