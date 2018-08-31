package actions

import (
	"github.com/4ydx/cdp/protocol/network"
	"github.com/4ydx/chrome-protocol"
	"log"
	"time"
)

// Cookies gets the browser's current cookies.
func Cookies(frame *cdp.Frame, timeout time.Duration) ([]network.Cookie, error) {
	action := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: network.CommandNetworkGetCookies, Params: &network.GetAllCookiesArgs{}, Reply: &network.GetAllCookiesReply{}, Timeout: timeout},
		})
	err := action.Run(frame)
	if err != nil {
		log.Print(err)
	}
	return action.Commands[0].Reply.(*network.GetAllCookiesReply).Cookies, err
}
