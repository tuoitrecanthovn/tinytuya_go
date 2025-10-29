package contrib

import (
	"encoding/json"
	"fmt"
	"time"

	"tinytuya_go/core"
)

const (
	IR_CMD_SEND_KEY_CODE = "send_ir"
	IR_DP_SEND_IR        = "201"
	IR_DP_LEARNED_ID     = "202"
)

// IRRemoteControlDevice represents a Tuya WiFi smart universal remote control simulator.
type IRRemoteControlDevice struct {
	*core.Device
	controlType int
}

// NewIRRemoteControlDevice creates a new IRRemoteControlDevice.
func NewIRRemoteControlDevice(d *core.Device) *IRRemoteControlDevice {
	return &IRRemoteControlDevice{Device: d}
}

// SendCommand sends a command to the device.
func (d *IRRemoteControlDevice) SendCommand(mode string, data map[string]interface{}) (map[string]interface{}, error) {
	if mode == "send" {
		var command map[string]interface{}
		if d.controlType == 1 {
			command = map[string]interface{}{
				"control": "send_ir",
				"type":    0,
			}
			if base64Code, ok := data["base64_code"]; ok {
				command["head"] = ""
				command["key1"] = "1" + base64Code.(string)
			} else if head, ok := data["head"]; ok {
				command["head"] = head
				command["key1"] = "0" + data["key"].(string)
			}
			jsonData, _ := json.Marshal(command)
			return d.SetValue(201, string(jsonData))
		} else if d.controlType == 2 {
			// Not implemented
			return nil, fmt.Errorf("controlType 2 not implemented")
		}
	} else if d.controlType == 1 {
		command := map[string]interface{}{"control": mode}
		jsonData, _ := json.Marshal(command)
		return d.SetValue(201, string(jsonData))
	} else if d.controlType == 2 {
		// Not implemented
		return nil, fmt.Errorf("controlType 2 not implemented")
	}
	return nil, fmt.Errorf("invalid mode or controlType")
}

// StudyStart starts a study session.
func (d *IRRemoteControlDevice) StudyStart() (map[string]interface{}, error) {
	return d.SendCommand("study", nil)
}

// StudyEnd ends a study session.
func (d *IRRemoteControlDevice) StudyEnd() (map[string]interface{}, error) {
	return d.SendCommand("study_exit", nil)
}


// SendButton simulates a learned button press.
func (d *IRRemoteControlDevice) SendButton(base64Code string) (map[string]interface{}, error) {
	return d.SendCommand("send", map[string]interface{}{"base64_code": base64Code})
}
