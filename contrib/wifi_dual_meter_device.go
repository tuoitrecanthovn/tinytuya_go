package contrib

import (
	"tinytuya_go/core"
)

const (
	WIFI_DUAL_METER_DPS_FORWARD_ENERGY_TOTAL = "1"
	WIFI_DUAL_METER_DPS_REVERSE_ENERGY_TOTAL = "2"
	WIFI_DUAL_METER_DPS_POWER_A              = "101"
	WIFI_DUAL_METER_DPS_POWER_B              = "105"
	WIFI_DUAL_METER_DPS_VOLTAGE              = "112"
	WIFI_DUAL_METER_DPS_CURRENT_A            = "113"
	WIFI_DUAL_METER_DPS_CURRENT_B            = "114"
	WIFI_DUAL_METER_DPS_TOTAL_POWER          = "115"
)

// WiFiDualMeterDevice represents a Tuya WiFi Dual Meter Device.
type WiFiDualMeterDevice struct {
	*core.Device
}

// GetValue returns a value from the device.
func (d *WiFiDualMeterDevice) GetValue(dpsCode string) (float64, error) {
	status, err := d.Status()
	if err != nil {
		return 0, err
	}
	val, _ := status["dps"].(map[string]interface{})[dpsCode].(float64)
	return val, nil
}

// GetForwardEnergyTotal returns the total forward energy.
func (d *WiFiDualMeterDevice) GetForwardEnergyTotal() (float64, error) {
	val, err := d.GetValue(WIFI_DUAL_METER_DPS_FORWARD_ENERGY_TOTAL)
	if err != nil {
		return 0, err
	}
	return val / 100, nil
}

// GetReverseEnergyTotal returns the total reverse energy.
func (d *WiFiDualMeterDevice) GetReverseEnergyTotal() (float64, error) {
	val, err := d.GetValue(WIFI_DUAL_METER_DPS_REVERSE_ENERGY_TOTAL)
	if err != nil {
		return 0, err
	}
	return val / 100, nil
}

// GetPowerA returns the power of phase A.
func (d *WiFiDualMeterDevice) GetPowerA() (float64, error) {
	val, err := d.GetValue(WIFI_DUAL_METER_DPS_POWER_A)
	if err != nil {
		return 0, err
	}
	return val / 10, nil
}
