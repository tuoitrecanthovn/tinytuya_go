package core

import (
	"encoding/json"
	"fmt"
	"time"
)

// Device represents a Tuya device with higher-level functions.
type Device struct {
	*XenonDevice
}

// NewDevice creates a new Device.
func NewDevice(devID, address, localKey, devType string, connectionTimeout time.Duration, version float64, persist bool, cid string, parent *XenonDevice) (*Device, error) {
	xenon, err := NewXenonDevice(devID, address, localKey, devType, connectionTimeout, version, persist, cid, parent)
	if err != nil {
		return nil, err
	}
	return &Device{XenonDevice: xenon}, nil
}

// SetStatus sets the status of the device to 'on' or 'off'.
func (d *Device) SetStatus(on bool, switchNum int) (map[string]interface{}, error) {
	payload, command := d.generatePayload(CONTROL, map[string]interface{}{
		fmt.Sprintf("%d", switchNum): on,
	})
	msg := TuyaMessage{
		Seqno:   d.seqno,
		Cmd:     uint32(command),
		Payload: payload,
	}
	d.seqno++
	data, err := d.sendReceive(msg)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// TurnOn turns the device on.
func (d *Device) TurnOn(switchNum int) (map[string]interface{}, error) {
	return d.SetStatus(true, switchNum)
}

// TurnOff turns the device off.
func (d *Device) TurnOff(switchNum int) (map[string]interface{}, error) {
	return d.SetStatus(false, switchNum)
}

// SetValue sets an integer value of any index.
func (d *Device) SetValue(index int, value interface{}) (map[string]interface{}, error) {
	payload, command := d.generatePayload(CONTROL, map[string]interface{}{
		fmt.Sprintf("%d", index): value,
	})
	msg := TuyaMessage{
		Seqno:   d.seqno,
		Cmd:     uint32(command),
		Payload: payload,
	}
	d.seqno++
	data, err := d.sendReceive(msg)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}
