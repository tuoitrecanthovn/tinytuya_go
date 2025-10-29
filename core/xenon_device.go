package core

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// XenonDevice represents a Tuya device.
type XenonDevice struct {
	ID                   string
	Address              string
	LocalKey             []byte
	DevType              string
	ConnectionTimeout    time.Duration
	Version              float64
	persist              bool
	CID                  string
	parent               *XenonDevice
	children             map[string]*XenonDevice
	port                 int
	socket               net.Conn
	socketPersistent     bool
	socketNODELAY        bool
	socketRetryLimit     int
	socketRetryDelay     time.Duration
	seqno                uint32
	dpsToRequest         map[string]interface{}
	autoIP               bool
	payloadDict          map[int]map[string]interface{}
}

// NewXenonDevice creates a new XenonDevice.
func NewXenonDevice(devID, address, localKey, devType string, connectionTimeout time.Duration, version float64, persist bool, cid string, parent *XenonDevice) (*XenonDevice, error) {
	d := &XenonDevice{
		ID:                devID,
		Address:           address,
		LocalKey:          []byte(localKey),
		DevType:           devType,
		ConnectionTimeout: connectionTimeout,
		Version:           version,
		persist:           persist,
		CID:               cid,
		parent:            parent,
		children:          make(map[string]*XenonDevice),
		port:              TCPPORT,
		socketPersistent:  persist,
		socketNODELAY:     true,
		socketRetryLimit:  5,
		socketRetryDelay:  5 * time.Second,
		seqno:             1,
		dpsToRequest:      make(map[string]interface{}),
	}

	if d.Address == "" || d.Address == "Auto" || d.Address == "0.0.0.0" {
		// Auto-discover IP address
		d.autoIP = true
		deviceInfo, err := FindDevice(devID)
		if err != nil {
			return nil, err
		}
		d.Address = deviceInfo["ip"].(string)
		// d.Version = deviceInfo["version"].(float64) // This will be a string, need to convert
	}

	return d, nil
}


// Status returns the device status.
func (d *XenonDevice) Status() (map[string]interface{}, error) {
	payload, command := d.generatePayload(DP_QUERY, nil)
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
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

var payloadDict = map[string]map[int]map[string]interface{}{
	"default": {
		CONTROL: {"command": map[string]interface{}{"devId": "", "uid": "", "t": ""}},
		STATUS:  {"command": map[string]interface{}{"gwId": "", "devId": ""}},
		DP_QUERY: {"command": map[string]interface{}{"gwId": "", "devId": "", "uid": "", "t": ""}},
	},
	"device22": {
		DP_QUERY: {
			"command_override": CONTROL_NEW,
			"command":          map[string]interface{}{"devId": "", "uid": "", "t": ""},
		},
	},
}

func (d *XenonDevice) generatePayload(command int, data map[string]interface{}) ([]byte, int) {
	jsonCommand, ok := payloadDict[d.DevType][command]["command"].(map[string]interface{})
	if !ok {
		return nil, 0
	}

	commandOverride, ok := payloadDict[d.DevType][command]["command_override"].(int)
	if !ok {
		commandOverride = command
	}

	jsonData := make(map[string]interface{})
	for k, v := range jsonCommand {
		jsonData[k] = v
	}

	if gwID, ok := jsonData["gwId"]; ok && gwID == "" {
		jsonData["gwId"] = d.ID
	}
	if devID, ok := jsonData["devId"]; ok && devID == "" {
		jsonData["devId"] = d.ID
	}
	if uid, ok := jsonData["uid"]; ok && uid == "" {
		jsonData["uid"] = d.ID
	}
	if t, ok := jsonData["t"]; ok && t == "" {
		jsonData["t"] = fmt.Sprintf("%d", time.Now().Unix())
	}

	if data != nil {
		jsonData["dps"] = data
	}

	payload, err := json.Marshal(jsonData)
	if err != nil {
		return nil, 0
	}

	return payload, commandOverride
}

func (d *XenonDevice) connect() error {
	if d.socket != nil {
		return nil
	}

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", d.Address, d.port), d.ConnectionTimeout)
	if err != nil {
		return err
	}
	d.socket = conn
	return nil
}

func (d *XenonDevice) sendReceive(msg TuyaMessage) ([]byte, error) {
	if err := d.connect(); err != nil {
		return nil, err
	}

	packed, err := PackMessage(msg, d.LocalKey)
	if err != nil {
		return nil, err
	}

	_, err = d.socket.Write(packed)
	if err != nil {
		return nil, err
	}

	response := make([]byte, 1024)
	n, err := d.socket.Read(response)
	if err != nil {
		return nil, err
	}

	unpacked, err := UnpackMessage(response[:n], d.LocalKey)
	if err != nil {
		return nil, err
	}

	return unpacked.Payload, nil
}
