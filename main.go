package main

import (
	"fmt"
	"log"
	"time"

	"tinytuya_go/core"
)

func main() {
	fmt.Println("=== TinyTuya Go Device Compatibility Test ===")

	const (
		DeviceID  = "ebaf701f96b67923b5vna9" // REPLACE WITH YOUR DEVICE ID
		IPAddress = "192.168.1.119"          // REPLACE WITH YOUR DEVICE IP
		LocalKey  = "Cq|@r(278+lDpzah"       // REPLACE WITH YOUR LOCAL KEY
	)

	// Validate inputs
	if len(LocalKey) != 16 {
		log.Printf("WARNING: Local key should be 16 characters for v3.4+, got %d characters", len(LocalKey))
	}

	// Try different protocol versions in order of complexity
	versions := []float64{3.3, 3.4, 3.5}
	var workingDevice *core.XenonDevice
	var workingVersion float64

	for _, version := range versions {
		fmt.Printf("\n=== Testing Protocol v%.1f ===\n", version)

		d, err := core.NewXenonDevice(
			DeviceID,
			IPAddress,
			LocalKey,
			"default", // devType
			10*time.Second, // Increased timeout
			version,
			false, // persist
			"",    // cid
			nil,   // parent
		)

		if err != nil {
			log.Printf("❌ Failed to create device (v%.1f): %v", version, err)
			continue
		}

		fmt.Printf("✅ Device created successfully with v%.1f\n", version)
		fmt.Printf("📡 Getting status with v%.1f...\n", version)

		data, err := d.Status()
		if err != nil {
			log.Printf("❌ Failed to get status (v%.1f): %v", version, err)

			// Try to close connection cleanly
			if d != nil {
				d.Close()
			}
			continue
		}

		fmt.Printf("🎉 SUCCESS with v%.1f!\n", version)
		fmt.Printf("📊 Device status: %+v\n", data)

		workingDevice = d
		workingVersion = version
		break
	}

	if workingDevice == nil {
		log.Fatalf("\n❌ All protocol versions failed. Please check:\n" +
			"   • Device IP address is correct\n" +
			"   • Device is powered on and connected to WiFi\n" +
			"   • Local key is correct (16 characters for v3.4+)\n" +
			"   • No firewall blocking TCP port 6668\n" +
			"   • Device and computer are on same network\n")
	}

	// If we have a working connection, try control operations
	fmt.Printf("\n=== Testing Control Operations with v%.1f ===\n", workingVersion)

	// Check if device has switch functionality
	status, err := workingDevice.Status()
	if err != nil {
		log.Printf("❌ Failed to get current status: %v", err)
		return
	}

	// Look for common switch DPS values
	switchDPS := findSwitchDPS(status)
	if switchDPS == "" {
		log.Println("⚠️  Could not find switch DPS, trying DPS '1' by default")
		switchDPS = "1"
	}

	currentState := getDPSValue(status, switchDPS)
	fmt.Printf("🔄 Current state of DPS %s: %v\n", switchDPS, currentState)

	// Toggle the switch
	var newState interface{}
	if currentState == true {
		newState = false
		fmt.Printf("🔴 Turning device OFF (DPS %s)...\n", switchDPS)
	} else {
		newState = true
		fmt.Printf("🟢 Turning device ON (DPS %s)...\n", switchDPS)
	}

	_, err = workingDevice.SetValue(switchDPS, newState)
	if err != nil {
		log.Printf("❌ Failed to set DPS %s to %v: %v", switchDPS, newState, err)
	} else {
		fmt.Printf("✅ Successfully set DPS %s to %v\n", switchDPS, newState)

		// Wait and verify
		time.Sleep(2 * time.Second)

		fmt.Println("🔍 Verifying new state...")
		newStatus, err := workingDevice.Status()
		if err != nil {
			log.Printf("❌ Failed to verify status: %v", err)
		} else {
			verifiedState := getDPSValue(newStatus, switchDPS)
			fmt.Printf("✅ Verified state of DPS %s: %v\n", switchDPS, verifiedState)

			if verifiedState == newState {
				fmt.Println("🎉 Control operation successful!")
			} else {
				fmt.Println("⚠️  State mismatch - device may take time to update")
			}
		}
	}

	// Test device22 mode if default failed
	if workingVersion == 3.2 || workingVersion == 3.3 {
		fmt.Println("\n=== Testing device22 mode ===")
		testDevice22Mode(DeviceID, IPAddress, LocalKey, workingVersion)
	}

	// Clean up
	workingDevice.Close()
	fmt.Println("\n✅ Test completed successfully!")
}

// findSwitchDPS attempts to find a switch DPS in the device status
func findSwitchDPS(status map[string]interface{}) string {
	if dps, ok := status["dps"].(map[string]interface{}); ok {
		// Common switch DPS values
		commonSwitchDPS := []string{"1", "20", "21", "22", "23", "24"}

		for _, dpsKey := range commonSwitchDPS {
			if _, exists := dps[dpsKey]; exists {
				fmt.Printf("🔍 Found potential switch DPS: %s\n", dpsKey)
				return dpsKey
			}
		}

		// If no common switches found, look for any boolean DPS
		for key, value := range dps {
			if _, isBool := value.(bool); isBool {
				fmt.Printf("🔍 Found boolean DPS: %s\n", key)
				return key
			}
		}
	}
	return ""
}

// getDPSValue safely extracts a DPS value from status
func getDPSValue(status map[string]interface{}, dpsKey string) interface{} {
	if dps, ok := status["dps"].(map[string]interface{}); ok {
		return dps[dpsKey]
	}
	return nil
}

// testDevice22Mode tries device22 configuration for devices that need it
func testDevice22Mode(deviceID, ipAddress, localKey string, version float64) {
	fmt.Println("🔄 Testing with device22 configuration...")

	d, err := core.NewXenonDevice(
		deviceID,
		ipAddress,
		localKey,
		"device22", // Use device22 type
		10*time.Second,
		version,
		false,
		"",
		nil,
	)

	if err != nil {
		log.Printf("❌ Failed to create device22: %v", err)
		return
	}
	defer d.Close()

	data, err := d.Status()
	if err != nil {
		log.Printf("❌ device22 mode failed: %v", err)
	} else {
		fmt.Printf("✅ device22 mode successful! Status: %+v\n", data)
	}
}
