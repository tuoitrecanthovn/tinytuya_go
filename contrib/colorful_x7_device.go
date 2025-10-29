package contrib

import (
	"fmt"

	"tinytuya_go/core"
)

const (
	COLORFULX7_DPS_INDEX_ON          = "20"
	COLORFULX7_DPS_INDEX_MODE        = "21"
	COLORFULX7_DPS_INDEX_COLOUR      = "24"
	COLORFULX7_DPS_INDEX_COUNTDOWN   = "26"
	COLORFULX7_DPS_INDEX_WORKMODE    = "104"
	COLORFULX7_DPS_INDEX_BRIGHTNESS  = "106"
	COLORFULX7_DPS_INDEX_DYNAMIC_MODE = "108"
	COLORFULX7_DPS_INDEX_MUSIC_MODE   = "109"
)

// ColorfulX7Device represents a Tuya based LED Music Controller.
type ColorfulX7Device struct {
	*core.Device
}

// IsOn returns the state of the device.
func (d *ColorfulX7Device) IsOn() (bool, error) {
	status, err := d.Status()
	if err != nil {
		return false, err
	}
	on, _ := status["dps"].(map[string]interface{})[COLORFULX7_DPS_INDEX_ON].(bool)
	return on, nil
}

// SwitchOff turns off the device.
func (d *ColorfulX7Device) SwitchOff() (map[string]interface{}, error) {
	return d.SetValue(20, false)
}

// SwitchOn turns on the device.
func (d *ColorfulX7Device) SwitchOn() (map[string]interface{}, error) {
	return d.SetValue(20, true)
}

// SetMode sets the mode to white | colour | scene | music | screen.
func (d *ColorfulX7Device) SetMode(mode string) (map[string]interface{}, error) {
	return d.SetValue(21, mode)
}

// SetColor sets the colour.
func (d *ColorfulX7Device) SetColor(r, g, b int) (map[string]interface{}, error) {
	// hsv := rgbToHSV(r, g, b)
	// hex := hsvToHex(hsv)
	// return d.SetValue(24, hex)
	return nil, fmt.Errorf("color conversion not implemented")
}

// SetBrightness sets the brightness.
func (d *ColorfulX7Device) SetBrightness(value int) (map[string]interface{}, error) {
	if value < 0 || value > 100 {
		return nil, fmt.Errorf("brightness must be between 0 and 100")
	}
	return d.SetValue(106, value)
}

// SetDynamicMode sets the scene type in DYNAMIC work mode.
func (d *ColorfulX7Device) SetDynamicMode(mode int) (map[string]interface{}, error) {
	if mode < 1 || mode > 180 {
		return nil, fmt.Errorf("dynamic mode must be between 1 and 180")
	}
	return d.SetValue(108, mode)
}

// SetMusicMode sets the scene type in MUSIC work mode.
func (d *ColorfulX7Device) SetMusicMode(mode int) (map[string]interface{}, error) {
	if mode < 1 || mode > 22 {
		return nil, fmt.Errorf("music mode must be between 1 and 22")
	}
	return d.SetValue(109, mode)
}

