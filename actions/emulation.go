package actions

import (
	"github.com/4ydx/cdp/protocol/emulation"
	"github.com/4ydx/chrome-protocol"
	"time"
)

func SetDeviceMetricsOverride(frame *cdp.Frame, width, height int, mobile bool, timeout time.Duration) error {
	// await client.Emulation.setDeviceMetricsOverride({width: 1920, height: 1080, fitWindow: true, deviceScaleFactor: 1, mobile: false});
	err := cdp.NewAction(
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: emulation.CommandEmulationSetDeviceMetricsOverride, Params: &emulation.SetDeviceMetricsOverrideArgs{
				Width:             width,
				Height:            height,
				DeviceScaleFactor: 1,
				Mobile:            mobile,
			}, Reply: &emulation.SetDeviceMetricsOverrideReply{}, Timeout: timeout},
		}).Run(frame)
	if err != nil {
		frame.Browser.Log.Print(err)
		return err
	}
	return nil
}
