package core

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
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
	sessionKey           []byte
	negotiatedSessionKey bool
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

// SetValue sets a single DPS value.
func (d *XenonDevice) SetValue(dpsID string, value interface{}) (map[string]interface{}, error) {
	payload, command := d.generatePayload(CONTROL, map[string]interface{}{dpsID: value})
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
		// It's common for control commands to return an empty or non-json payload
		return nil, nil
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

	if d.Version >= 3.4 {
		return d.negotiateSessionKey()
	}

	return nil
}

func (d *XenonDevice) negotiateSessionKey() error {
	// Step 1: Send client nonce
	clientNonce := make([]byte, 16)
	_, err := rand.Read(clientNonce)
	if err != nil {
		return fmt.Errorf("failed to generate client nonce: %w", err)
	}

	startMsg := TuyaMessage{
		Seqno:   d.seqno,
		Cmd:     SESS_KEY_NEG_START,
		Payload: clientNonce,
	}
	d.seqno++

	packedStart, err := PackPlaintext55AA(startMsg)
	if err != nil {
		return fmt.Errorf("failed to pack start message: %w", err)
	}

	_, err = d.socket.Write(packedStart)
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// Step 2: Receive device nonce and HMAC
	response := make([]byte, 1024)
	n, err := d.socket.Read(response)
	if err != nil {
		return fmt.Errorf("failed to read response to start message: %w", err)
	}

	unpackedResp, err := UnpackPlaintext55AA(response[:n])
	if err != nil {
		return fmt.Errorf("failed to unpack response message: %w", err)
	}

	if unpackedResp.Cmd != uint32(SESS_KEY_NEG_RESP) {
		return fmt.Errorf("unexpected command in response: got %d, want %d", unpackedResp.Cmd, SESS_KEY_NEG_RESP)
	}

	deviceNonce := unpackedResp.Payload[:16]
	hmacFromDevice := unpackedResp.Payload[16:]

	// Verify HMAC
	mac := hmac.New(sha256.New, d.LocalKey)
	mac.Write(clientNonce)
	expectedHMAC := mac.Sum(nil)

	if !hmac.Equal(hmacFromDevice, expectedHMAC) {
		return fmt.Errorf("HMAC verification failed")
	}

	// Step 3: Send HMAC of device nonce
	mac.Reset()
	mac.Write(deviceNonce)
	hmacToDevice := mac.Sum(nil)

	finishMsg := TuyaMessage{
		Seqno:   d.seqno,
		Cmd:     uint32(SESS_KEY_NEG_FINISH),
		Payload: hmacToDevice,
	}
	d.seqno++

	packedFinish, err := PackPlaintext55AA(finishMsg)
	if err != nil {
		return fmt.Errorf("failed to pack finish message: %w", err)
	}

	_, err = d.socket.Write(packedFinish)
	if err != nil {
		return fmt.Errorf("failed to send finish message: %w", err)
	}

	// Key Derivation for v3.5
	tmpKey := make([]byte, 16)
	for i := 0; i < 16; i++ {
		tmpKey[i] = deviceNonce[i] ^ clientNonce[i]
	}

	ciphertext, tag, err := GCMEncrypt(d.LocalKey, clientNonce[:12], tmpKey, nil)
	if err != nil {
		return fmt.Errorf("failed to derive session key: %w", err)
	}

	combined := append(ciphertext, tag...)
	d.sessionKey = combined[12:28]
	d.negotiatedSessionKey = true

	return nil
}

func (d *XenonDevice) sendReceive(msg TuyaMessage) ([]byte, error) {
	if err := d.connect(); err != nil {
		return nil, err
	}

	var packed []byte
	var err error

	if d.negotiatedSessionKey {
		packed, err = PackMessage6699(msg, d.sessionKey)
	} else {
		packed, err = PackMessage(msg, d.LocalKey)
	}

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

	var unpacked *TuyaMessage
	if d.negotiatedSessionKey {
		unpacked, err = UnpackMessage6699(response[:n], d.sessionKey)
	} else {
		unpacked, err = UnpackMessage(response[:n], d.LocalKey)
	}

	if err != nil {
		return nil, err
	}

	return unpacked.Payload, nil
}
