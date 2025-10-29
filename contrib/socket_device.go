package contrib

import (
	"tinytuya_go/core"
)

const (
	SOCKET_DPS_STATE   = "1"
	SOCKET_DPS_CURRENT = "18"
	SOCKET_DPS_POWER   = "19"
	SOCKET_DPS_VOLTAGE = "20"
)

// SocketDevice represents a Tuya based Socket.
type SocketDevice struct {
	*core.Device
}

// GetEnergyConsumption returns the energy consumption data.
func (d *SocketDevice) GetEnergyConsumption() (map[string]interface{}, error) {
	status, err := d.Status()
	if err != nil {
		return nil, err
	}
	return status["dps"].(map[string]interface{}), nil
}

// GetCurrent returns the current in mA.
func (d *SocketDevice) GetCurrent() (float64, error) {
	status, err := d.Status()
	if err != nil {
		return 0, err
	}
	current, _ := status["dps"].(map[string]interface{})[SOCKET_DPS_CURRENT].(float64)
	return current, nil
}

// GetPower returns the power in W.
func (d *SocketDevice) GetPower() (float64, error) {
	status, err := d.Status()
	if err != nil {
		return 0, err
	}
	power, _ := status["dps"].(map[string]interface{})[SOCKET_DPS_POWER].(float64)
	return power / 10, nil
}

// GetVoltage returns the voltage in V.
func (d *SocketDevice) GetVoltage() (float64, error) {
	status, err := d.Status()
	if err != nil {
		return 0, err
	}
	voltage, _ := status["dps"].(map[string]interface{})[SOCKET_DPS_VOLTAGE].(float64)
	return voltage / 10, nil
}

// GetState returns the current state of the device.
func (d *SocketDevice) GetState() (bool, error) {
	status, err := d.Status()
	if err != nil {
		return false, err
	}
	on, _ := status["dps"].(map[string]interface{})[SOCKET_DPS_STATE].(bool)
	return on, nil
}
