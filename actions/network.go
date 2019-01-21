package actions

import (
	"github.com/4ydx/cdp/protocol/network"
	"github.com/4ydx/chrome-protocol"
	"net/http"
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
		frame.Browser.Log.Print(err)
	}
	return action.Commands[0].Reply.(*network.GetAllCookiesReply).Cookies, err
}

// SetCookie sets one cookie in the browser.
func SetCookie(frame *cdp.Frame, url string, cookie *http.Cookie, timeout time.Duration) (bool, error) {
	tse := network.TimeSinceEpoch(float64(cookie.Expires.Unix()))
	params := &network.SetCookieArgs{
		URL:      url,
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		Secure:   cookie.Secure,
		HTTPOnly: cookie.HttpOnly,
		Expires:  &tse,
	}
	action := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: network.CommandNetworkSetCookie, Params: params, Reply: &network.SetCookieReply{}, Timeout: timeout},
		})
	err := action.Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
	}
	return action.Commands[0].Reply.(*network.SetCookieReply).Success, err
}
