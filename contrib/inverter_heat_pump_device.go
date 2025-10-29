package contrib

import (
	"tinytuya_go/core"
)

// TemperatureUnit represents the unit of the temperature.
type TemperatureUnit bool

const (
	CELSIUS    TemperatureUnit = true
	FAHRENHEIT TemperatureUnit = false
)

// InverterHeatPumpMode represents the mode of the inverter.
type InverterHeatPumpMode string

const (
	HEATING InverterHeatPumpMode = "warm"
	UNKNOWN InverterHeatPumpMode = "unknown"
)

// InverterHeatPumpFault represents the fault of the inverter.
type InverterHeatPumpFault int

const (
	NOMINAL       InverterHeatPumpFault = 0
	NO_WATER_FLOW InverterHeatPumpFault = 4
	UNKNOWN_FAULT InverterHeatPumpFault = -1
)

const (
	INVERTER_DPS_ON_DP                        = "1"
	INVERTER_DPS_INLET_WATER_TEMP_DP          = "102"
	INVERTER_DPS_UNIT_DP                        = "103"
	INVERTER_DPS_HEATING_CAPACITY_PERCENT_DP    = "104"
	INVERTER_DPS_MODE_DP                        = "105"
	INVERTER_DPS_TARGET_WATER_TEMP_DP         = "106"
	INVERTER_DPS_LOWER_LIMIT_TARGET_WATER_TEMP_DP = "107"
	INVERTER_DPS_UPPER_LIMIT_TARGET_WATER_TEMP_DP = "108"
	INVERTER_DPS_FAULT_DP                       = "115"
	INVERTER_DPS_SILENCE_MODE_DP                = "117"
)

// InverterHeatPumpDevice represents a Tuya WiFi smart inverter heat pump.
type InverterHeatPumpDevice struct {
	*core.Device
}

// IsOn returns True if the inverter is on.
func (d *InverterHeatPumpDevice) IsOn() (bool, error) {
	status, err := d.Status()
	if err != nil {
		return false, err
	}
	on, _ := status["dps"].(map[string]interface{})[INVERTER_DPS_ON_DP].(bool)
	return on, nil
}

// GetUnit returns the unit of the temperature.
func (d *InverterHeatPumpDevice) GetUnit() (TemperatureUnit, error) {
	status, err := d.Status()
	if err != nil {
		return false, err
	}
	unit, _ := status["dps"].(map[string]interface{})[INVERTER_DPS_UNIT_DP].(bool)
	return TemperatureUnit(unit), nil
}

// SetTargetWaterTemp sets the target water temperature.
func (d *InverterHeatPumpDevice) SetTargetWaterTemp(target float64) (map[string]interface{}, error) {
	return d.SetValue(106, target)
}
