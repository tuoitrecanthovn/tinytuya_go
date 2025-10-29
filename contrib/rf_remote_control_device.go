package contrib

import (
	"fmt"

	"tinytuya_go/core"
)

// RFRemoteControlDevice represents a Tuya WiFi smart universal RF remote controller.
type RFRemoteControlDevice struct {
	*IRRemoteControlDevice
}

// NewRFRemoteControlDevice creates a new RFRemoteControlDevice.
func NewRFRemoteControlDevice(d *core.Device) *RFRemoteControlDevice {
	ir := NewIRRemoteControlDevice(d)
	return &RFRemoteControlDevice{IRRemoteControlDevice: ir}
}

// RFStudyStart starts an RF study session.
func (d *RFRemoteControlDevice) RFStudyStart(freq int, short bool) (map[string]interface{}, error) {
	data := map[string]interface{}{"freq": fmt.Sprintf("%d", freq)}
	cmd := "rf_study"
	if short {
		cmd = "rf_shortstudy"
	}
	return d.SendCommand(cmd, data)
}

// RFStudyEnd ends an RF study session.
func (d *RFRemoteControlDevice) RFStudyEnd(freq int, short bool) (map[string]interface{}, error) {
	data := map[string]interface{}{"freq": fmt.Sprintf("%d", freq)}
	cmd := "rfstudy_exit"
	if short {
		cmd = "rfshortstudy_exit"
	}
	return d.SendCommand(cmd, data)
}

// RFSendButton sends a learned RF button press.
func (d *RFRemoteControlDevice) RFSendButton(base64Code string, times, delay, intervals int) (map[string]interface{}, error) {
	key1 := map[string]interface{}{
		"code":      base64Code,
		"times":     times,
		"delay":     delay,
		"intervals": intervals,
	}
	data := map[string]interface{}{"key1": key1}
	return d.SendCommand("rfstudy_send", data)
}
