package contrib

import (
	"fmt"

	"tinytuya_go/core"
)

const (
	DOORBELL_DPS_BASIC_INDICATOR   = "101"
	DOORBELL_DPS_VOLUME              = "160"
	DOORBELL_DPS_MOTION_AREA         = "169"
	DOORBELL_DPS_MOTION_AREA_SWITCH = "168"
)

// DoorbellDevice represents a Tuya based Video-Doorbell.
type DoorbellDevice struct {
	*core.Device
}

// SetBasicIndicator sets the basic indicator.
func (d *DoorbellDevice) SetBasicIndicator(val bool) (map[string]interface{}, error) {
	return d.SetValue(101, val)
}

// SetVolume sets the doorbell volume.
func (d *DoorbellDevice) SetVolume(vol int) (map[string]interface{}, error) {
	if vol < 3 {
		vol = 3
	}
	if vol > 10 {
		vol = 10
	}
	return d.SetValue(160, vol)
}

// SetMotionArea sets the area of motion detection.
func (d *DoorbellDevice) SetMotionArea(x, y, xlen, ylen int) (map[string]interface{}, error) {
	data := fmt.Sprintf(`{"num":1,"region0":{"x":%d,"y":%d,"xlen":%d,"ylen":%d}}`, x, y, xlen, ylen)
	return d.SetValue(169, data)
}

// SetMotionAreaSwitch turns the motion detection area on/off.
func (d *DoorbellDevice) SetMotionAreaSwitch(useArea bool) (map[string]interface{}, error) {
	return d.SetValue(168, useArea)
}
