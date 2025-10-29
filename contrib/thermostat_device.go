package contrib

import (
	"tinytuya_go/core"
)

const (
	THERMOSTAT_DPS_MODE             = "2"
	THERMOSTAT_DPS_TEMP_SET         = "16"
	THERMOSTAT_DPS_TEMP_CURRENT     = "24"
	THERMOSTAT_DPS_UPPER_TEMP       = "108"
	THERMOSTAT_DPS_LOWER_TEMP       = "109"
	THERMOSTAT_DPS_FAN              = "115"
	THERMOSTAT_DPS_SCHEDULE_ENABLED = "119"
	THERMOSTAT_DPS_HOLD             = "120"
)

// ThermostatDevice represents a Tuya based 24v Thermostat.
type ThermostatDevice struct {
	*core.Device
}

// SetSetpoint sets the target temperature.
func (d *ThermostatDevice) SetSetpoint(setpoint float64) (map[string]interface{}, error) {
	return d.SetValue(16, setpoint)
}

// SetCoolSetpoint sets the cooling setpoint.
func (d *ThermostatDevice) SetCoolSetpoint(setpoint float64) (map[string]interface{}, error) {
	return d.SetValue(108, setpoint)
}

// SetHeatSetpoint sets the heating setpoint.
func (d *ThermostatDevice) SetHeatSetpoint(setpoint float64) (map[string]interface{}, error) {
	return d.SetValue(109, setpoint)
}

// SetMode sets the system mode.
func (d *ThermostatDevice) SetMode(mode string) (map[string]interface{}, error) {
	return d.SetValue(2, mode)
}

// SetFan sets the fan mode.
func (d *ThermostatDevice) SetFan(fan string) (map[string]interface{}, error) {
	return d.SetValue(115, fan)
}

// SetSchedule enables or disables the schedule.
func (d *ThermostatDevice) SetSchedule(enabled bool) (map[string]interface{}, error) {
	return d.SetValue(119, enabled)
}

// SetHold sets the temperature hold.
func (d *ThermostatDevice) SetHold(hold string) (map[string]interface{}, error) {
	return d.SetValue(120, hold)
}
