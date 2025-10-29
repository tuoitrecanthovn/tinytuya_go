package contrib

import (
	"tinytuya_go/core"
)

const (
	DPS_MODE         = "101"
	DPS_CUR_TEMP     = "102"
	DPS_SWITCH_STATE = "103"
	DPS_CURRENT      = "108"
	DPS_POWER        = "109"
	DPS_VOLTAGE      = "110"
	DPS_TEMP_UNIT    = "118"
	DPS_TOTAL_POWER  = "111"
)

// AtorchTemperatureControllerDevice represents a Tuya ATORCH-Temperature Controller.
type AtorchTemperatureControllerDevice struct {
	*core.Device
}

// GetEnergyConsumption returns the energy consumption data.
func (d *AtorchTemperatureControllerDevice) GetEnergyConsumption() (map[string]interface{}, error) {
	status, err := d.Status()
	if err != nil {
		return nil, err
	}
	return status["dps"].(map[string]interface{}), nil
}

// GetCurrent returns the current in mA.
func (d *AtorchTemperatureControllerDevice) GetCurrent() (float64, error) {
	status, err := d.Status()
	if err != nil {
		return 0, err
	}
	current, _ := status["dps"].(map[string]interface{})[DPS_CURRENT].(float64)
	return current, nil
}

// GetPower returns the power in W.
func (d *AtorchTemperatureControllerDevice) GetPower() (float64, error) {
	status, err := d.Status()
	if err != nil {
		return 0, err
	}
	power, _ := status["dps"].(map[string]interface{})[DPS_POWER].(float64)
	return power / 100, nil
}

// GetVoltage returns the voltage in V.
func (d *AtorchTemperatureControllerDevice) GetVoltage() (float64, error) {
	status, err := d.Status()
	if err != nil {
		return 0, err
	}
	voltage, _ := status["dps"].(map[string]interface{})[DPS_VOLTAGE].(float64)
	return voltage / 100, nil
}

// GetTemp returns the current temperature.
func (d *AtorchTemperatureControllerDevice) GetTemp() (float64, error) {
	status, err := d.Status()
	if err != nil {
		return 0, err
	}
	temp, _ := status["dps"].(map[string]interface{})[DPS_CUR_TEMP].(float64)
	return temp / 10, nil
}

// GetState returns the current state of the device.
func (d *AtorchTemperatureControllerDevice) GetState() (map[string]interface{}, error) {
	status, err := d.Status()
	if err != nil {
		return nil, err
	}
	mode, _ := status["dps"].(map[string]interface{})[DPS_MODE].(string)
	if mode == "socket" {
		state, _ := status["dps"].(map[string]interface{})[DPS_SWITCH_STATE].(bool)
		statusStr := "off"
		if state {
			statusStr = "on"
		}
		return map[string]interface{}{"mode": mode, "status": statusStr}, nil
	}
	return map[string]interface{}{"mode": mode}, nil
}
