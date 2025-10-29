package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// FindDevice scans the network for a Tuya device with a specific ID.
func FindDevice(devID string) (map[string]interface{}, error) {
	devices, err := DeviceScan(false, 3)
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if device["gwId"] == devID {
			return device, nil
		}
	}

	return nil, fmt.Errorf("device not found")
}

// DeviceScan scans the network for Tuya devices.
func DeviceScan(verbose bool, maxRetry int) (map[string]map[string]interface{}, error) {
	// This is a simplified implementation of the UDP scanner.
	// A full implementation would require more complex logic to handle different subnets and network interfaces.

	addr, err := net.ResolveUDPAddr("udp", ":"+fmt.Sprintf("%d", UDPPORT))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	devices := make(map[string]map[string]interface{})

	// Send broadcast packets
	broadcastAddr := &net.UDPAddr{IP: net.IPv4bcast, Port: UDPPORT}
	payload := []byte("{\"gwId\":\"\",\"t\":\"0\"}") // Simplified payload

	for i := 0; i < maxRetry; i++ {
		conn.WriteToUDP(payload, broadcastAddr)
	}

	// Listen for responses
	conn.SetReadDeadline(time.Now().Add(time.Duration(maxRetry) * time.Second))
	buffer := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			// Timeout reached
			break
		}

		decrypted, err := DecryptUDP(buffer[:n])
		if err != nil {
			continue
		}

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(decrypted), &result); err == nil {
			if ip, ok := result["ip"].(string); ok {
				devices[ip] = result
			}
		}
	}

	return devices, nil
}

// DecryptUDP decrypts a UDP message from a Tuya device.
func DecryptUDP(msg []byte) (string, error) {
	// UDP broadcasts are not always encrypted
	if len(msg) > 0 && msg[0] == '{' {
		return string(msg), nil
	}

	// 3.1, 3.3, and 3.4 UDP broadcasts are encrypted with a static key
	key := []byte{0x79, 0x47, 0x41, 0x64, 0x6c, 0x6f, 0x70, 0x6f, 0x50, 0x56, 0x6c, 0x64, 0x41, 0x42, 0x66, 0x6e} // md5(yGAdlopoPVldABfn)

	// Check for 3.3+ header
	if bytes.HasPrefix(msg, []byte(PROTOCOL_VERSION_BYTES_33)) || bytes.HasPrefix(msg, []byte(PROTOCOL_VERSION_BYTES_34)) {
		// Strip header and decrypt
		msg = msg[len(PROTOCOL_VERSION_BYTES_33)+len(PROTOCOL_3x_HEADER):]
		decrypted, err := ECBDecrypt(key, msg)
		if err != nil {
			return "", err
		}
		return string(decrypted), nil
	}

	// 3.1 UDP broadcasts are encrypted without a header
	decrypted, err := ECBDecrypt(key, msg)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}
