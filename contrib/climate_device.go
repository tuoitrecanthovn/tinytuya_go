package contrib

import (
	"fmt"

	"tinytuya_go/core"
)

const (
	CLIMATE_DPS_POWER     = "1"
	CLIMATE_DPS_SET_TEMP  = "2"
	CLIMATE_DPS_CUR_TEMP  = "3"
	CLIMATE_DPS_MODE      = "4"
	CLIMATE_DPS_FAN       = "5"
	CLIMATE_DPS_TEMP_UNIT = "19"
	CLIMATE_DPS_TIMER     = "22"
	CLIMATE_DPS_STATE     = "101"
)

// ClimateDevice represents a Tuya based Air Conditioner.
type ClimateDevice struct {
	*core.Device
}

// GetRoomTemperature returns the room temperature.
func (d *ClimateDevice) GetRoomTemperature() (float64, error) {
	status, err := d.Status()
	if err != nil {
		return 0, err
	}
	temp, _ := status["dps"].(map[string]interface{})[CLIMATE_DPS_CUR_TEMP].(float64)
	return temp, nil
}

// GetTargetTemperature returns the target temperature.
func (d *ClimateDevice) GetTargetTemperature() (float64, error) {
	status, err := d.Status()
	if err != nil {
		return 0, err
	}
	temp, _ := status["dps"].(map[string]interface{})[CLIMATE_DPS_SET_TEMP].(float64)
	return temp, nil
}

// SetTargetTemperature sets the target temperature.
func (d *ClimateDevice) SetTargetTemperature(t float64) (map[string]interface{}, error) {
	return d.SetValue(2, t)
}

// GetOperatingMode returns the operating mode.
func (d *ClimateDevice) GetOperatingMode() (string, error) {
	status, err := d.Status()
	if err != nil {
		return "", err
	}
	mode, _ := status["dps"].(map[string]interface{})[CLIMATE_DPS_MODE].(string)
	return mode, nil
}

// SetOperatingMode sets the operating mode.
func (d *ClimateDevice) SetOperatingMode(mode string) (map[string]interface{}, error) {
	if mode != "cold" && mode != "hot" && mode != "dehumidify" {
		return nil, fmt.Errorf("invalid mode")
	}
	return d.SetValue(4, mode)
}
