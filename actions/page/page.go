package pa

import (
	"github.com/4ydx/cdproto"
	"github.com/4ydx/cdproto/page"
	"github.com/4ydx/chrome-protocol"
	"time"
)

// Navigate sends the browser to the given URL
func Navigate(id *cdp.ID, url string, timeout time.Duration) *cdp.Action {
	return cdp.NewAction(
		[]cdp.Event{
			cdp.Event{Name: cdproto.EventPageFrameStartedLoading, Value: &page.EventFrameStartedLoading{}, IsRequired: true},
			cdp.Event{Name: cdproto.EventPageFrameStoppedLoading, Value: &page.EventFrameStoppedLoading{}, IsRequired: true},
		},
		[]cdp.Step{
			cdp.Step{Id: id.GetNext(), Method: page.CommandNavigate, Params: &page.NavigateParams{URL: "https://google.com"}, Returns: &page.NavigateReturns{}, Timeout: timeout},
		})
}
