package page

import (
	"github.com/4ydx/cdp/protocol/page"
	"github.com/4ydx/chrome-protocol"
	"time"
)

// Navigate sends the browser to the given URL
func Navigate(pg *cdp.Page, url string, timeout time.Duration) error {
	return cdp.NewAction(pg,
		[]cdp.Event{
			cdp.Event{Name: page.EventPageFrameStartedLoading, Value: &page.FrameStartedLoadingReply{}, IsRequired: true},
			cdp.Event{Name: page.EventPageFrameStoppedLoading, Value: &page.FrameStoppedLoadingReply{}, IsRequired: true},
		},
		[]cdp.Step{
			cdp.Step{ID: pg.ID.GetNext(), Method: page.CommandPageNavigate, Params: &page.NavigateArgs{URL: url}, Reply: &page.NavigateReply{}, Timeout: timeout},
		}).Run()
}
